package integration

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HDIOES/cpa-backend/models"

	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/pkg/errors"
)

const (
	shikimoriError          = "Connection to shikimori error"
	readBytesError          = "Read bytes error"
	unmarshalError          = "Json unmarshal error"
	parseErrorMessage       = "Parse error"
	animeUpdateErrorMessage = "Anime update error"
)

//ShikimoriJob struct
type ShikimoriJob struct {
	AnimeDao  *models.AnimeDAO
	GenreDao  *models.GenreDAO
	StudioDao *models.StudioDAO
	Config    *util.Configuration
}

//Run function
func (sj *ShikimoriJob) Run() {
	defer log.Println("Job has been ended")
	client := &http.Client{}
	//at start we need to load studios and genres
	if processStudioErr := sj.ProcessStudios(client); processStudioErr != nil {
		util.HandleError(processStudioErr)
		return
	}
	time.Sleep(1000 * time.Millisecond)
	if processGenresErr := sj.ProcessGenres(client); processGenresErr != nil {
		util.HandleError(processGenresErr)
		return
	}
	time.Sleep(1000 * time.Millisecond)
	//then we have to load anime list
	animes := []Anime{}
	var page int64 = 1
	for len(animes) == 50 || page == 1 {
		animesPatch, animesErr := sj.ProcessAnimePatch(page, client)
		if animesErr != nil {
			util.HandleError(animesErr)
		}
		animes = animesPatch
		page++
		time.Sleep(1000 * time.Millisecond)
	}
	//then we need to run long loading of animes by call url '/api/animes/:id'
	var animesDtos, err = sj.GetNotProcessedExternalAnimes()
	if err != nil {
		util.HandleError(err)
		return
	}
	for _, animeDto := range animesDtos {
		processOneAmineErr := sj.ProcessOneAnime(client, animeDto)
		if processOneAmineErr != nil {
			util.HandleError(processOneAmineErr)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

//GetNotProcessedExternalAnimes function
func (sj *ShikimoriJob) GetNotProcessedExternalAnimes() ([]models.AnimeDTO, error) {
	sqlBuilder := models.AnimeQueryBuilder{}
	sqlBuilder.SetProcessed(false)
	animeDtos, getAnimeDtosErr := sj.AnimeDao.FindByFilter(sqlBuilder)
	if getAnimeDtosErr != nil {
		return nil, errors.Wrap(getAnimeDtosErr, "FindByFilter error")
	}
	return animeDtos, nil
}

//ProcessOneAnime function
func (sj *ShikimoriJob) ProcessOneAnime(client *http.Client, animeDto models.AnimeDTO) error {
	log.Println("Now we will process anime with external_id = " + animeDto.ExternalID)
	resp, getAnimeByIDErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriAnimeSearchURL + "/" + animeDto.ExternalID)
	if getAnimeByIDErr != nil {
		return errors.Wrap(getAnimeByIDErr, shikimoriError)
	}
	defer resp.Body.Close()
	anime := Anime{}
	body, readStudiosErr := ioutil.ReadAll(resp.Body)
	if readStudiosErr != nil {
		return errors.Wrap(readStudiosErr, readBytesError)
	}
	log.Println("Response body: ", string(body))
	parseError := json.Unmarshal(body, &anime)
	if parseError != nil {
		return errors.Wrap(parseError, unmarshalError)
	}
	animeDto.Rating = anime.Rating
	if anime.Score != nil {
		score, parseErr := strconv.ParseFloat(*anime.Score, 64)
		if parseErr == nil {
			animeDto.Score = &score
		} else {
			util.HandleError(errors.Wrap(parseErr, parseErrorMessage))
		}
	}
	animeDto.Franchise = anime.Franchise
	if anime.Duration != nil {
		durationFloat := float64(*anime.Duration)
		animeDto.Duration = &durationFloat
	}
	processed := true
	animeDto.Processed = &processed
	//then we need to update row in database
	updateErr := sj.AnimeDao.Update(animeDto)
	if updateErr != nil {
		return errors.Wrap(updateErr, animeUpdateErrorMessage)
	}
	//and now let go to set genre for anime
	for _, g := range *anime.Genres {
		genreDto, genreDtoErr := sj.GenreDao.FindByExternalID(strconv.FormatInt(*g.ID, 10))
		if genreDtoErr != nil {
			return errors.Wrap(genreDtoErr, "")
		}
		if linkErr := sj.AnimeDao.LinkAnimeAndGenre(animeDto.ID, genreDto.ID); linkErr != nil {
			return errors.Wrap(linkErr, "")
		}
	}
	//let go to set studio for anime
	for _, s := range *anime.Studios {
		studioDto, studioDtoErr := sj.StudioDao.FindByExternalID(strconv.FormatInt(*s.ID, 10))
		if studioDtoErr != nil {
			return errors.Wrap(studioDtoErr, "")
		}
		if linkErr := sj.AnimeDao.LinkAnimeAndStudio(animeDto.ID, studioDto.ID); linkErr != nil {
			return errors.Wrap(linkErr, "")
		}
	}
	log.Println("Anime has been processed")
	return nil
}

//ProcessGenres function
func (sj *ShikimoriJob) ProcessGenres(client *http.Client) error {
	genres := []Genre{}
	resp, getGenresErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriGenreURL)
	if getGenresErr != nil {
		return errors.Wrap(getGenresErr, "")
	}
	defer resp.Body.Close()
	body, readGenresErr := ioutil.ReadAll(resp.Body)
	if readGenresErr != nil {
		return errors.Wrap(readGenresErr, "")
	}
	parseGenresError := json.Unmarshal(body, &genres)
	if parseGenresError != nil {
		return errors.Wrap(parseGenresError, "")
	}
	for _, genre := range genres {
		externalID := strconv.FormatInt(*genre.ID, 10)
		genreDto, dtoErr := sj.GenreDao.FindByExternalID(externalID)
		genreNotFound := true
		if genreNotFound {
			genreNotFound = strings.Compare(dtoErr.Error(), "Genre not found") == 0
		}
		dto := models.GenreDTO{}
		dto.ExternalID = externalID
		dto.Name = genre.Name
		dto.Russian = genre.Russian
		dto.Kind = genre.Kind
		if genreNotFound {
			_, createErr := sj.GenreDao.Create(dto)
			if createErr != nil {
				return errors.Wrap(createErr, "")
			}
		} else {
			dto.ID = genreDto.ID
			updateErr := sj.GenreDao.Update(dto)
			if updateErr != nil {
				return errors.Wrap(updateErr, "")
			}
		}
	}
	log.Println("Genres have been processed")
	return nil
}

//ProcessStudios function
func (sj *ShikimoriJob) ProcessStudios(client *http.Client) error {
	studios := []Studio{}
	resp, getStudioErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriStudioURL)
	if getStudioErr != nil {
		return errors.Wrap(getStudioErr, shikimoriError)
	}
	defer resp.Body.Close()
	body, readStudiosErr := ioutil.ReadAll(resp.Body)
	if readStudiosErr != nil {
		return errors.Wrap(getStudioErr, readBytesError)
	}
	parseError := json.Unmarshal(body, &studios)
	if parseError != nil {
		return errors.Wrap(parseError, unmarshalError)
	}
	for _, shikiStudio := range studios {
		externalID := strconv.FormatInt(*shikiStudio.ID, 10)
		studioDto, findErr := sj.StudioDao.FindByExternalID(externalID)
		studioNotFound := true
		if findErr != nil {
			studioNotFound = strings.Compare(findErr.Error(), "Studio not found") == 0
		}
		dto := models.StudioDTO{
			ExternalID:         externalID,
			Name:               shikiStudio.Name,
			FilteredStudioName: shikiStudio.FilteredName,
			IsReal:             shikiStudio.Real,
			ImageURL:           shikiStudio.Image,
		}
		if studioNotFound {
			if _, createErr := sj.StudioDao.Create(dto); createErr != nil {
				return errors.Wrap(createErr, "")
			}
		} else {
			dto.ID = studioDto.ID
			if updateErr := sj.StudioDao.Update(dto); updateErr != nil {
				return errors.Wrap(updateErr, "")
			}
		}
	}
	log.Println("Studios have been processed")
	return nil
}

//ProcessAnimePatch function
func (sj *ShikimoriJob) ProcessAnimePatch(page int64, client *http.Client) ([]Anime, error) {
	animes := []Anime{}
	resp, animesGetErr := client.Get(sj.Config.ShikimoriURL + sj.Config.ShikimoriAnimeSearchURL + "?page=" + strconv.FormatInt(page, 10) + "&limit=50")
	if animesGetErr != nil {
		return nil, errors.Wrap(animesGetErr, "")
	}
	defer resp.Body.Close()
	body, readAnimesErr := ioutil.ReadAll(resp.Body)
	if readAnimesErr != nil {
		return nil, errors.Wrap(readAnimesErr, "")
	}
	parseAnimesError := json.Unmarshal(body, &animes)
	if parseAnimesError != nil {
		return nil, errors.Wrap(parseAnimesError, "")
	}
	for _, anime := range animes {
		animeDto, animeDtoErr := sj.AnimeDao.FindByExternalID(strconv.FormatInt(*anime.ID, 10))
		dto := models.AnimeDTO{}
		dto.ExternalID = strconv.FormatInt(*anime.ID, 10)
		dto.Name = anime.Name
		dto.Russian = anime.Russian
		dto.Kind = anime.Kind
		dto.PosterURL = anime.Image.Original
		dto.AnimeURL = anime.URL
		dto.Kind = anime.Kind
		dto.Status = anime.Status
		dto.Epizodes = anime.Episodes
		dto.EpizodesAired = anime.EpisodesAired

		airedOn := anime.AiredOn.Local()
		dto.AiredOn = &airedOn

		releasedOn := anime.ReleasedOn.Local()
		dto.ReleasedOn = &releasedOn

		processed := false
		dto.Processed = &processed
		animeNotFound := true
		if animeDtoErr != nil {
			animeNotFound = strings.Compare(animeDtoErr.Error(), "Anime not found") == 0
		}
		if animeNotFound {
			if _, createErr := sj.AnimeDao.Create(dto); createErr != nil {
				return nil, errors.Wrap(createErr, "")
			}
		} else {
			dto.ID = animeDto.ID
			if updateErr := sj.AnimeDao.Update(dto); updateErr != nil {
				return nil, errors.Wrap(updateErr, "")
			}
		}
	}
	log.Println("Page with number " + strconv.FormatInt(page, 10) + " has been processed")
	return animes, nil
}

func rollbackTransaction(tx *sql.Tx, err error) error {
	if rollbackErr := tx.Rollback(); err != nil {
		return errors.Wrap(rollbackErr, "")
	}
	return errors.Wrap(err, "")
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
	return errors.Wrap(err, "")
}

func (sts *ShikimoriTime) toDateValue() *string {
	value := sts.Format("2006-01-02")
	return &value
}
