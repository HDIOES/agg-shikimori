package integration

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/pkg/errors"
)

//ShikimoriDao struct
type ShikimoriDao struct {
	Client *http.Client
	Config *util.Configuration
}

//Animes func
func (sh *ShikimoriDao) Animes(page, limit int) ([]Anime, error) {
	animes := []Anime{}
	resp, animesGetErr := sh.Client.Get(sh.Config.ShikimoriURL + sh.Config.ShikimoriAnimeSearchURL + "?page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit))
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
	return animes, nil
}

//OneAnime func
func (sh *ShikimoriDao) OneAnime(id string) (*Anime, error) {
	resp, getAnimeByIDErr := sh.Client.Get(sh.Config.ShikimoriURL + sh.Config.ShikimoriAnimeSearchURL + "/" + id)
	if getAnimeByIDErr != nil {
		return nil, errors.Wrap(getAnimeByIDErr, "")
	}
	defer resp.Body.Close()
	anime := Anime{}
	body, readAnimesErr := ioutil.ReadAll(resp.Body)
	if readAnimesErr != nil {
		return nil, errors.Wrap(readAnimesErr, "")
	}
	parseError := json.Unmarshal(body, &anime)
	if parseError != nil {
		return nil, errors.Wrap(parseError, unmarshalError)
	}
	return &anime, nil
}

//Studios func
func (sh *ShikimoriDao) Studios() ([]Studio, error) {
	studios := []Studio{}
	resp, getStudioErr := sh.Client.Get(sh.Config.ShikimoriURL + sh.Config.ShikimoriStudioURL)
	if getStudioErr != nil {
		return nil, errors.Wrap(getStudioErr, "")
	}
	defer resp.Body.Close()
	body, readStudioErr := ioutil.ReadAll(resp.Body)
	if readStudioErr != nil {
		return nil, errors.Wrap(readStudioErr, "")
	}
	parseStudiosError := json.Unmarshal(body, &studios)
	if parseStudiosError != nil {
		return nil, errors.Wrap(parseStudiosError, "")
	}
	return studios, nil
}

//Genres func
func (sh *ShikimoriDao) Genres() ([]Genre, error) {
	genres := []Genre{}
	resp, getGenreErr := sh.Client.Get(sh.Config.ShikimoriURL + sh.Config.ShikimoriGenreURL)
	if getGenreErr != nil {
		return nil, errors.Wrap(getGenreErr, "")
	}
	defer resp.Body.Close()
	body, readGenreErr := ioutil.ReadAll(resp.Body)
	if readGenreErr != nil {
		return nil, errors.Wrap(readGenreErr, "")
	}
	parseGenresError := json.Unmarshal(body, &genres)
	if parseGenresError != nil {
		return nil, errors.Wrap(parseGenresError, "")
	}
	return genres, nil
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
