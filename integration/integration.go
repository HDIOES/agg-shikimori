package integration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/HDIOES/cpa-backend/util"
)

var NDD = errors.New("Database does not contains rows with processed = false") //'No database data' error

type ShikimoriJob struct {
	Db     *sql.DB
	Config util.Configuration
}

func (sj *ShikimoriJob) Run() {
	client := &http.Client{}
	//at start we need to load studios and genres
	sj.ProcessStudios(client)
	time.Sleep(1000 * time.Millisecond)
	sj.ProcessGenres(client)
	time.Sleep(1000 * time.Millisecond)
	//then we have to load anime list
	animes := &[]Anime{}
	var page int64 = 1
	for len(*animes) == 50 || page == 1 {
		animes = sj.ProcessAnimePatch(page, client)
		page++
		time.Sleep(1000 * time.Millisecond)
	}
	//then we need to run long loading of animes by call url '/api/animes/:id'
	var externalAnimeIDs, err = sj.GetNotProcessedExternalAnimeIds()
	if err != nil {
		log.Println("Error getting of anime ids: ", err)
		return
	}
	for _, eID := range *externalAnimeIDs {
		processOneAmineErr := sj.ProcessOneAnime(client, eID)
		if processOneAmineErr != nil {
			log.Println("Error getting of anime: ", processOneAmineErr)
		}
		time.Sleep(1000 * time.Millisecond)
	}
	log.Println("Job has been ended")
}

func (sj *ShikimoriJob) GetNotProcessedExternalAnimeIds() (externalAnimeIDs *[]string, err error) {
	getAnimeIdsRows, getAnimeIdsErr := sj.Db.Query("SELECT external_id FROM anime WHERE processed = false")
	if getAnimeIdsErr != nil {
		log.Println("Error getting of anime ids: ", getAnimeIdsErr)
		return nil, getAnimeIdsErr
	}
	defer getAnimeIdsRows.Close()
	var ids []string
	var externlalID sql.NullString
	for getAnimeIdsRows.Next() {
		getAnimeIdsRows.Scan(&externlalID)
		ids = append(ids, externlalID.String)
	}
	return &ids, nil
}

//ProcessOneAnime function
func (sj *ShikimoriJob) ProcessOneAnime(client *http.Client, eID string) (err error) {
	tx, txErr := sj.Db.Begin()
	if txErr != nil {
		log.Println("Transaction start failed: ", txErr)
		err = txErr
		return
	}
	defer func(tx *sql.Tx) {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}(tx)

	log.Println("Now we will process anime with external_id = " + eID)
	resp, getAnimeByIDErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriAnimeSearchURL + eID)
	if getAnimeByIDErr != nil {
		log.Println("Error during getting anime by id: ", getAnimeByIDErr)
		err = getAnimeByIDErr
		panic(getAnimeByIDErr)
	}
	defer resp.Body.Close()
	anime := &Anime{}
	body, readStudiosErr := ioutil.ReadAll(resp.Body)
	if readStudiosErr != nil {
		log.Println("Error during reading studios: ", readStudiosErr)
		err = readStudiosErr
		panic(readStudiosErr)
	}
	log.Println("Response body: ", string(body))
	parseError := json.Unmarshal(body, anime)
	if parseError != nil {
		log.Println("Error during parsing anime: ", parseError)
		err = parseError
		panic(readStudiosErr)
	}
	//then we need to update row in database
	var score *float64
	floatScore, parseScoreErr := strconv.ParseFloat(*(anime.Score), 32)
	if parseScoreErr != nil {
		log.Println("Error during parsing score: ", parseScoreErr)
		score = nil
	} else {
		score = &floatScore
	}
	_, execTxErr := tx.Exec("UPDATE anime SET score = $1, duration = $2, rating = $3, franchase = $4, processed = true, lastmodifytime = now() WHERE external_id = $5",
		score, anime.Duration, anime.Rating, anime.Franchise, eID)
	if execTxErr != nil {
		log.Println("Query cannot be executed: ", execTxErr)
		err = execTxErr
		panic(execTxErr)
	}
	//and now let go to set genre for anime
	for _, g := range *(anime.Genres) {
		animeGenreRows, findGenreErr := tx.Query("SELECT anime_genre.anime_id, anime_genre.genre_id "+
			"FROM anime_genre "+
			"JOIN anime ON anime.id = anime_genre.anime_id "+
			"JOIN genre ON genre.id = anime_genre.genre_id "+
			"WHERE anime.external_id = $1 AND genre.external_id = $2",
			strconv.FormatInt(*(anime.ID), 10),
			strconv.FormatInt(*(g.ID), 10))
		if findGenreErr != nil {
			log.Println("Query cannot be executed: ", findGenreErr)
			err = findGenreErr
			panic(findGenreErr)
		}
		if !animeGenreRows.Next() {
			animeGenreRows.Close()
			//now we need insert missing genre
			_, insertNewGenreForAnime := tx.Exec("INSERT INTO anime_genre (anime_id, genre_id) "+
				"SELECT anime.id, genre.id FROM anime JOIN genre ON anime.external_id = $1 AND genre.external_id = $2",
				strconv.FormatInt(*(anime.ID), 10),
				strconv.FormatInt(*(g.ID), 10))
			if insertNewGenreForAnime != nil {
				log.Println("Query cannot be executed: ", insertNewGenreForAnime)
				err = insertNewGenreForAnime
				panic(insertNewGenreForAnime)
			}
		} else {
			animeGenreRows.Close()
		}
	}
	//let go to set studio for anime
	for _, s := range *(anime.Studios) {
		animeStudioRows, findStudioErr := tx.Query("SELECT anime_studio.anime_id, anime_studio.studio_id FROM anime_studio "+
			"join anime on anime.id = anime_studio.anime_id join studio on studio.id = anime_studio.studio_id WHERE anime.external_id = $1 AND studio.external_id = $2",
			strconv.FormatInt(*(anime.ID), 10),
			strconv.FormatInt(*(s.ID), 10))
		if findStudioErr != nil {
			log.Println("Query cannot be executed: ", findStudioErr)
			err = findStudioErr
			panic(findStudioErr)
		}
		if !animeStudioRows.Next() {
			animeStudioRows.Close()
			//now we need insert missing studio
			_, insertNewStudioForAnime := tx.Exec("INSERT INTO anime_studio SELECT anime.id, studio.id FROM anime JOIN studio ON anime.external_id = $1 AND studio.external_id = $2",
				strconv.FormatInt(*(anime.ID), 10),
				strconv.FormatInt(*(s.ID), 10))
			if insertNewStudioForAnime != nil {
				log.Println("Query cannot be executed: ", insertNewStudioForAnime)
				err = insertNewStudioForAnime
				panic(insertNewStudioForAnime)
			}
		} else {
			animeStudioRows.Close()
		}
	}

	if txCommitErr := tx.Commit(); txCommitErr != nil {
		log.Println("Transaction cannot be commited: ", txCommitErr)
		err = txCommitErr
		panic(txCommitErr)
	}
	log.Println("Anime has been processed")
	return
}

//ProcessGenres function
func (sj *ShikimoriJob) ProcessGenres(client *http.Client) {
	tx, txErr := sj.Db.Begin()
	if txErr != nil {
		log.Println("Transaction start failed: ", txErr)
		return
	}
	defer func(tx *sql.Tx) {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}(tx)
	genres := &[]Genre{}
	resp, getGenresErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriGenreURL)
	if getGenresErr != nil {
		log.Println("Error during getting genres: ", getGenresErr)
		panic(getGenresErr)
	}
	defer resp.Body.Close()
	body, readGenresErr := ioutil.ReadAll(resp.Body)
	if readGenresErr != nil {
		log.Println("Error during reading genres: ", readGenresErr)
		panic(readGenresErr)
	}
	parseGenresError := json.Unmarshal(body, genres)
	if parseGenresError != nil {
		log.Println("Error during parsing genres: ", parseGenresError)
		panic(parseGenresError)
	}
	for i := 0; i < len(*genres); i++ {
		rows, txExecSelectErr := tx.Query("SELECT external_id FROM genre WHERE external_id = $1", (*genres)[i].ID)
		if txExecSelectErr != nil {
			log.Println("Query cannot be executed: ", txExecSelectErr)
			panic(parseGenresError)
		}
		if !rows.Next() {
			_, txExecErr := tx.Exec("INSERT INTO genre (external_id, genre_name, russian, kind) "+
				"VALUES ($1, $2, $3, $4)",
				(*genres)[i].ID,
				(*genres)[i].Name,
				(*genres)[i].Russian,
				(*genres)[i].Kind)
			if txExecErr != nil {
				log.Println("Query cannot be executed: ", txExecErr)
				rows.Close()
				panic(txExecErr)
			}
		}
		rows.Close()
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		log.Println("Transaction cannot be commited: ", txCommitErr)
		panic(txCommitErr)
	}
	log.Println("Genres have been processed")
}

//ProcessStudios function
func (sj *ShikimoriJob) ProcessStudios(client *http.Client) {
	tx, txErr := sj.Db.Begin()
	if txErr != nil {
		log.Println("Transaction start failed: ", txErr)
		return
	}
	defer func(tx *sql.Tx) {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}(tx)
	studios := &[]Studio{}
	resp, getStudioErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriStudioURL)
	if getStudioErr != nil {
		log.Println("Error during getting studios: ", getStudioErr)
		panic(getStudioErr)
	}
	defer resp.Body.Close()
	body, readStudiosErr := ioutil.ReadAll(resp.Body)
	if readStudiosErr != nil {
		log.Println("Error during reading studios: ", readStudiosErr)
		panic(readStudiosErr)
	}
	parseError := json.Unmarshal(body, studios)
	if parseError != nil {
		log.Println("Error during parsing studios: ", parseError)
		panic(readStudiosErr)
	}
	for i := 0; i < len(*studios); i++ {
		rows, txExecSelectErr := tx.Query("SELECT external_id FROM studio WHERE external_id = $1", (*studios)[i].ID)
		if txExecSelectErr != nil {
			log.Println("Query cannot be executed: ", txExecSelectErr)
			panic(readStudiosErr)
		}
		if !rows.Next() {
			_, txExecErr := tx.Exec("INSERT INTO studio (external_id, studio_name, filtered_studio_name, is_real, image_url) "+
				"VALUES ($1, $2, $3, $4, $5)",
				(*studios)[i].ID,
				(*studios)[i].Name,
				(*studios)[i].FilteredName,
				(*studios)[i].Real,
				(*studios)[i].Image)
			if txExecErr != nil {
				log.Println("Query cannot be executed: ", txExecErr)
				rows.Close()
				panic(txExecErr)
			}
		}
		rows.Close()
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		log.Println("Transaction cannot be commited: ", txCommitErr)
		panic(txCommitErr)
	}
	log.Println("Studios have been processed")
}

//ProcessAnimePatch function
func (sj *ShikimoriJob) ProcessAnimePatch(page int64, client *http.Client) *[]Anime {
	animes := &[]Anime{}
	tx, txErr := sj.Db.Begin()
	if txErr != nil {
		log.Println("Transaction start failed: ", txErr)
		return animes
	}
	defer func(tx *sql.Tx) {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}(tx)

	resp, animesGetErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriAnimeSearchURL + "?page=" + strconv.FormatInt(page, 10) + "&limit=50")
	if animesGetErr != nil {
		log.Println("Error during getting animes: ", animesGetErr)
		panic(animesGetErr)
	}
	defer resp.Body.Close()
	body, readAnimesErr := ioutil.ReadAll(resp.Body)
	if readAnimesErr != nil {
		log.Println("Error during reading body: ", readAnimesErr)
		panic(readAnimesErr)
	}
	parseAnimesError := json.Unmarshal(body, animes)
	if parseAnimesError != nil {
		log.Println("Error parsing of animes: ", parseAnimesError)
		panic(parseAnimesError)
	}
	//function for inserting anime
	insertAnimeFunc := func(tx *sql.Tx, anime Anime) {
		rows, txExecSelectErr := tx.Query("SELECT external_id FROM ANIME WHERE external_id = $1", anime.ID)
		if txExecSelectErr != nil {
			log.Println("Query cannot be executed: ", txExecSelectErr)
			panic(txExecSelectErr)
		}
		defer rows.Close()
		if !rows.Next() {
			var airedOn *string
			if anime.AiredOn != nil {
				airedOn = anime.AiredOn.toDateValue()
			}
			var releasedOn *string
			if anime.ReleasedOn != nil {
				releasedOn = anime.ReleasedOn.toDateValue()
			}
			var posterURL = *(anime.Image.Original)
			_, txExecErr := tx.Exec("INSERT INTO anime (external_id, name, russian, amine_url, kind, anime_status, epizodes, epizodes_aired, aired_on, released_on, poster_url, processed, lastmodifytime) "+
				"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, false, now())",
				anime.ID,
				anime.Name,
				anime.Russian,
				anime.URL,
				anime.Kind,
				anime.Status,
				anime.Episodes,
				anime.EpisodesAired,
				airedOn,
				releasedOn,
				posterURL)
			if txExecErr != nil {
				log.Println("Query cannot be executed: ", txExecErr)
				panic(txExecErr)
			}
		}
	}
	for i := 0; i < len(*animes); i++ {
		insertAnimeFunc(tx, (*animes)[i])
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		log.Println("Transaction cannot be commited: ", txCommitErr)
		panic(txCommitErr)
	}
	log.Println("Page with number " + strconv.FormatInt(page, 10) + " has been processed")
	return animes
}

//Anime struct
type Anime struct {
	ID                 *int64                          `json:"id"`
	Name               *string                         `json:"name"`
	Russian            *string                         `json:"russian"`
	Image              *Image                          `json:"image"`
	URL                *string                         `json:"url"`
	Kind               *string                         `json:"kind"`
	Status             *string                         `json:"status"`
	Episodes           *int64                          `json:"episodes"`
	EpisodesAired      *int64                          `json:"episodes_aired"`
	AiredOn            *ShikimoriTime                  `json:"aired_on"`
	ReleasedOn         *ShikimoriTime                  `json:"released_on"`
	Rating             *string                         `json:"rating"`
	English            *[]string                       `json:"english"`
	Japanese           *[]string                       `json:"japanese"`
	Synonyms           *[]string                       `json:"synonyms"`
	LicenseNameRu      *string                         `json:"license_name_ru"`
	Duration           *int64                          `json:"duration"`
	Score              *string                         `json:"score"`
	Description        *string                         `json:"description"`
	DescriptionHTML    *string                         `json:"description_html"`
	DescriptionSource  *string                         `json:"description_source"`
	Franchise          *string                         `json:"franchise"`
	Favoured           *bool                           `json:"favoured"`
	Anons              *bool                           `json:"anons"`
	Ongoing            *bool                           `json:"ongoing"`
	ThreadID           *int64                          `json:"thread_id"`
	TopicID            *int64                          `json:"topic_id"`
	MyAnimelistID      *int64                          `json:"myanimelist_id"`
	RatesScoresStats   *[]RatesScoresStatsNameValue    `json:"rates_scores_stats"`
	RatesStatusesStats *[]RatesScoresStatusesNameValue `json:"rates_statuses_stats"`
	UpdatedAt          *string                         `json:"updated_at"`      //NEED TO CHANGE ON DATETIME!!!!
	NextEpisodeAt      *string                         `json:"next_episode_at"` //NEED TO CHANGE ON DATETIME!!!!
	Genres             *[]Genre                        `json:"genres"`
	Studios            *[]Studio                       `json:"studios"`
	Videos             *[]Video                        `json:"videos"`
	Screenshots        *[]Screenshot                   `json:"screenshots"`
	Userrate           *string                         `json:"userrate"`
}

//RatesScoresStatsNameValue struct
type RatesScoresStatsNameValue struct {
	Name  *int64 `json:"name"`
	Value *int64 `json:"value"`
}

//RatesScoresStatusesNameValue struct
type RatesScoresStatusesNameValue struct {
	Name  *string `json:"name"`
	Value *int64  `json:"value"`
}

//Studio struct
type Studio struct {
	ID           *int64  `json:"id"`
	Name         *string `json:"name"`
	FilteredName *string `json:"filtered_name"`
	Real         *bool   `json:"real"`
	Image        *string `json:"image"`
}

//Genre struct
type Genre struct {
	ID      *int64  `json:"id"`
	Name    *string `json:"name"`
	Russian *string `json:"russian"`
	Kind    *string `json:"kind"`
}

//Video struct
type Video struct {
	ID        *int64  `json:"id"`
	URL       *string `json:"url"`
	ImageURL  *string `json:"image_url"`
	PlayerURL *string `json:"player_url"`
	Name      *string `json:"name"`
	Kind      *string `json:"kind"`
	Hosting   *string `json:"hosting"`
}

//Screenshot struct
type Screenshot struct {
	Original *string `json:"original"`
	Preview  *string `json:"preview"`
}

//Image struct
type Image struct {
	Original *string `json:"original"`
	Preview  *string `json:"preview"`
	X96      *string `json:"x96"`
	X48      *string `json:"x48"`
}

//ShikimoriTime struct
type ShikimoriTime struct {
	time.Time
}

//UnmarshalJSON unmarshales ShikimoriTime correctly
func (sts *ShikimoriTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	s = s[1 : len(s)-1]

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.999999999Z0700", s)
	}
	sts.Time = t
	return err
}

func (sts *ShikimoriTime) toDateValue() *string {
	value := sts.Format("2006-01-02")
	return &value
}
