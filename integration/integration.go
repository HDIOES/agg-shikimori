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

	"github.com/HDIOES/cpa-backend/rest/util"
)

var NDD = errors.New("Database does not contains rows with processed = false") //'No database data' error

type ShikimoriJob struct {
	Db     *sql.DB
	Config *util.Configuration
}

//Run function
func (sj *ShikimoriJob) Run() {
	defer log.Println("Job has been ended")
	client := &http.Client{}
	//at start we need to load studios and genres
	if processStudioErr := sj.ProcessStudios(client); processStudioErr != nil {
		log.Print("Studios processing error", processStudioErr)
		return
	}
	time.Sleep(1000 * time.Millisecond)
	if processGenresErr := sj.ProcessGenres(client); processGenresErr != nil {
		log.Print("Genres processing error", processGenresErr)
		return
	}
	time.Sleep(1000 * time.Millisecond)
	//then we have to load anime list
	animes := &[]Anime{}
	var page int64 = 1
	for len(*animes) == 50 || page == 1 {
		animesPatch, animesErr := sj.ProcessAnimePatch(page, client)
		if animesErr != nil {
			log.Print("Error anime patch processing", animesErr)
		}
		animes = animesPatch
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
}

//GetNotProcessedExternalAnimeIds function
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
func (sj *ShikimoriJob) ProcessOneAnime(client *http.Client, eID string) error {
	tx, txErr := sj.Db.Begin()
	if txErr != nil {
		return rollbackTransaction(tx, txErr)
	}
	log.Println("Now we will process anime with external_id = " + eID)
	resp, getAnimeByIDErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriAnimeSearchURL + "/" + eID)
	if getAnimeByIDErr != nil {
		return rollbackTransaction(tx, getAnimeByIDErr)
	}
	defer resp.Body.Close()
	anime := &Anime{}
	body, readStudiosErr := ioutil.ReadAll(resp.Body)
	if readStudiosErr != nil {
		return rollbackTransaction(tx, readStudiosErr)
	}
	log.Println("Response body: ", string(body))
	parseError := json.Unmarshal(body, anime)
	if parseError != nil {
		return rollbackTransaction(tx, parseError)
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
		return rollbackTransaction(tx, execTxErr)
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
			return rollbackTransaction(tx, findGenreErr)
		}
		if !animeGenreRows.Next() {
			animeGenreRows.Close()
			//now we need insert missing genre
			_, insertNewGenreForAnime := tx.Exec("INSERT INTO anime_genre (anime_id, genre_id) "+
				"SELECT anime.id, genre.id FROM anime JOIN genre ON anime.external_id = $1 AND genre.external_id = $2",
				strconv.FormatInt(*(anime.ID), 10),
				strconv.FormatInt(*(g.ID), 10))
			if insertNewGenreForAnime != nil {
				return rollbackTransaction(tx, insertNewGenreForAnime)
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
			return rollbackTransaction(tx, findStudioErr)
		}
		if !animeStudioRows.Next() {
			animeStudioRows.Close()
			//now we need insert missing studio
			_, insertNewStudioForAnime := tx.Exec("INSERT INTO anime_studio SELECT anime.id, studio.id FROM anime JOIN studio ON anime.external_id = $1 AND studio.external_id = $2",
				strconv.FormatInt(*(anime.ID), 10),
				strconv.FormatInt(*(s.ID), 10))
			if insertNewStudioForAnime != nil {
				return rollbackTransaction(tx, insertNewStudioForAnime)
			}
		} else {
			animeStudioRows.Close()
		}
	}

	if txCommitErr := tx.Commit(); txCommitErr != nil {
		return rollbackTransaction(tx, txCommitErr)
	}
	log.Println("Anime has been processed")
	return nil
}

//ProcessGenres function
func (sj *ShikimoriJob) ProcessGenres(client *http.Client) error {
	tx, txErr := sj.Db.Begin()
	if txErr != nil {
		return rollbackTransaction(tx, txErr)
	}
	genres := &[]Genre{}
	resp, getGenresErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriGenreURL)
	if getGenresErr != nil {
		return rollbackTransaction(tx, getGenresErr)
	}
	defer resp.Body.Close()
	body, readGenresErr := ioutil.ReadAll(resp.Body)
	if readGenresErr != nil {
		return rollbackTransaction(tx, readGenresErr)
	}
	parseGenresError := json.Unmarshal(body, genres)
	if parseGenresError != nil {
		return rollbackTransaction(tx, parseGenresError)
	}
	for i := 0; i < len(*genres); i++ {
		rows, txExecSelectErr := tx.Query("SELECT external_id FROM genre WHERE external_id = $1", (*genres)[i].ID)
		if txExecSelectErr != nil {
			return rollbackTransaction(tx, txExecSelectErr)
		}
		if !rows.Next() {
			_, txExecErr := tx.Exec("INSERT INTO genre (external_id, genre_name, russian, kind) "+
				"VALUES ($1, $2, $3, $4)",
				(*genres)[i].ID,
				(*genres)[i].Name,
				(*genres)[i].Russian,
				(*genres)[i].Kind)
			if txExecErr != nil {
				rows.Close()
				return rollbackTransaction(tx, txExecErr)
			}
		}
		rows.Close()
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		return rollbackTransaction(tx, txCommitErr)
	}
	log.Println("Genres have been processed")
	return nil
}

//ProcessStudios function
func (sj *ShikimoriJob) ProcessStudios(client *http.Client) error {
	studios := &[]Studio{}
	tx, txErr := sj.Db.Begin()
	if txErr != nil {
		return rollbackTransaction(tx, txErr)
	}
	resp, getStudioErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriStudioURL)
	if getStudioErr != nil {
		return rollbackTransaction(tx, getStudioErr)
	}
	defer resp.Body.Close()
	body, readStudiosErr := ioutil.ReadAll(resp.Body)
	if readStudiosErr != nil {
		return rollbackTransaction(tx, readStudiosErr)
	}
	parseError := json.Unmarshal(body, studios)
	if parseError != nil {
		return rollbackTransaction(tx, parseError)
	}
	for i := 0; i < len(*studios); i++ {
		rows, txExecSelectErr := tx.Query("SELECT external_id FROM studio WHERE external_id = $1", (*studios)[i].ID)
		if txExecSelectErr != nil {
			return rollbackTransaction(tx, txExecSelectErr)
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
				rows.Close()
				return rollbackTransaction(tx, txExecErr)
			}
		}
		rows.Close()
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		return rollbackTransaction(tx, txCommitErr)
	}
	log.Println("Studios have been processed")
	return nil
}

//ProcessAnimePatch function
func (sj *ShikimoriJob) ProcessAnimePatch(page int64, client *http.Client) (*[]Anime, error) {
	animes := &[]Anime{}
	tx, txErr := sj.Db.Begin()
	if txErr != nil {
		return nil, rollbackTransaction(tx, txErr)
	}
	resp, animesGetErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriAnimeSearchURL + "?page=" + strconv.FormatInt(page, 10) + "&limit=50")
	if animesGetErr != nil {
		return nil, rollbackTransaction(tx, animesGetErr)
	}
	defer resp.Body.Close()
	body, readAnimesErr := ioutil.ReadAll(resp.Body)
	if readAnimesErr != nil {
		return nil, rollbackTransaction(tx, readAnimesErr)
	}
	parseAnimesError := json.Unmarshal(body, animes)
	if parseAnimesError != nil {
		return nil, rollbackTransaction(tx, parseAnimesError)
	}
	//function for inserting anime
	insertAnimeFunc := func(tx *sql.Tx, anime Anime) error {
		rows, txExecSelectErr := tx.Query("SELECT external_id FROM ANIME WHERE external_id = $1", anime.ID)
		if txExecSelectErr != nil {
			return txExecSelectErr
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
				return txExecErr
			}
		}
		return nil
	}
	for i := 0; i < len(*animes); i++ {
		if err := insertAnimeFunc(tx, (*animes)[i]); err != nil {
			return nil, rollbackTransaction(tx, err)
		}
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		return nil, rollbackTransaction(tx, txCommitErr)
	}
	log.Println("Page with number " + strconv.FormatInt(page, 10) + " has been processed")
	return animes, nil
}

func rollbackTransaction(tx *sql.Tx, err error) error {
	if rollbackErr := tx.Rollback(); err != nil {
		return rollbackErr
	}
	return err
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
