package integration

import (
	"database/sql"
	"log"
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
	AnimeDao     *models.AnimeDAO
	GenreDao     *models.GenreDAO
	StudioDao    *models.StudioDAO
	Config       *util.Configuration
	ShikimoriDao *ShikimoriDao
}

//Run function
func (sj *ShikimoriJob) Run() {
	defer log.Println("Job has been ended")
	//at start we need to load studios and genres
	if processStudioErr := sj.ProcessStudios(); processStudioErr != nil {
		util.HandleError(processStudioErr)
		return
	}
	time.Sleep(1000 * time.Millisecond)
	if processGenresErr := sj.ProcessGenres(); processGenresErr != nil {
		util.HandleError(processGenresErr)
		return
	}
	time.Sleep(1000 * time.Millisecond)
	//then we have to load anime list
	animes := []Anime{}
	page := 1
	for len(animes) == 50 || page == 1 {
		animesPatch, animesErr := sj.ProcessAnimePatch(page)
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
		processOneAmineErr := sj.ProcessOneAnime(animeDto)
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
func (sj *ShikimoriJob) ProcessOneAnime(animeDto models.AnimeDTO) error {
	anime, err := sj.ShikimoriDao.OneAnime(animeDto.ExternalID)
	if err != nil {
		return errors.Wrap(err, "")
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
func (sj *ShikimoriJob) ProcessGenres() error {
	genres, err := sj.ShikimoriDao.Genres()
	if err != nil {
		return errors.Wrap(err, "")
	}
	for _, genre := range genres {
		externalID := strconv.FormatInt(*genre.ID, 10)
		genreDto, dtoErr := sj.GenreDao.FindByExternalID(externalID)
		genreNotFound := false
		if dtoErr != nil {
			genreNotFound = strings.Compare(dtoErr.Error(), "Genre not found") == 0
		}
		if genreNotFound {
			dto := models.GenreDTO{}
			dto.ExternalID = externalID
			dto.Name = genre.Name
			dto.Russian = genre.Russian
			dto.Kind = genre.Kind
			_, createErr := sj.GenreDao.Create(dto)
			if createErr != nil {
				return errors.Wrap(createErr, "")
			}
		} else {
			genreDto.ExternalID = externalID
			genreDto.Name = genre.Name
			genreDto.Russian = genre.Russian
			genreDto.Kind = genre.Kind
			updateErr := sj.GenreDao.Update(*genreDto)
			if updateErr != nil {
				return errors.Wrap(updateErr, "")
			}
		}
	}
	log.Println("Genres have been processed")
	return nil
}

//ProcessStudios function
func (sj *ShikimoriJob) ProcessStudios() error {
	studios, err := sj.ShikimoriDao.Studios()
	if err != nil {
		return errors.Wrap(err, "")
	}
	for _, shikiStudio := range studios {
		externalID := strconv.FormatInt(*shikiStudio.ID, 10)
		studioDto, findErr := sj.StudioDao.FindByExternalID(externalID)
		studioNotFound := false
		if findErr != nil {
			studioNotFound = strings.Compare(findErr.Error(), "Studio not found") == 0
		}
		if studioNotFound {
			dto := models.StudioDTO{
				ExternalID:         externalID,
				Name:               shikiStudio.Name,
				FilteredStudioName: shikiStudio.FilteredName,
				IsReal:             shikiStudio.Real,
				ImageURL:           shikiStudio.Image,
			}
			if _, createErr := sj.StudioDao.Create(dto); createErr != nil {
				return errors.Wrap(createErr, "")
			}
		} else {
			studioDto.ExternalID = strconv.FormatInt(*shikiStudio.ID, 10)
			studioDto.Name = shikiStudio.Name
			studioDto.FilteredStudioName = shikiStudio.FilteredName
			studioDto.IsReal = shikiStudio.Real
			studioDto.ImageURL = shikiStudio.Image
			if updateErr := sj.StudioDao.Update(*studioDto); updateErr != nil {
				return errors.Wrap(updateErr, "")
			}
		}
	}
	log.Println("Studios have been processed")
	return nil
}

//ProcessAnimePatch function
func (sj *ShikimoriJob) ProcessAnimePatch(page int) ([]Anime, error) {
	animes, err := sj.ShikimoriDao.Animes(page, 50)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	for _, anime := range animes {
		animeDto, animeDtoErr := sj.AnimeDao.FindByExternalID(strconv.FormatInt(*anime.ID, 10))
		animeNotFound := false
		if animeDtoErr != nil {
			animeNotFound = strings.Compare(animeDtoErr.Error(), "Anime not found") == 0
		}
		if animeNotFound {
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
			if anime.AiredOn != nil {
				airedOn := anime.AiredOn.Local()
				dto.AiredOn = &airedOn
			}
			if anime.ReleasedOn != nil {
				releasedOn := anime.ReleasedOn.Local()
				dto.ReleasedOn = &releasedOn
			}
			processed := false
			dto.Processed = &processed
			if _, createErr := sj.AnimeDao.Create(dto); createErr != nil {
				return nil, errors.Wrap(createErr, "")
			}
		} else {
			animeDto.ExternalID = strconv.FormatInt(*anime.ID, 10)
			animeDto.Name = anime.Name
			animeDto.Russian = anime.Russian
			animeDto.Kind = anime.Kind
			animeDto.PosterURL = anime.Image.Original
			animeDto.AnimeURL = anime.URL
			animeDto.Kind = anime.Kind
			animeDto.Status = anime.Status
			animeDto.Epizodes = anime.Episodes
			animeDto.EpizodesAired = anime.EpisodesAired
			if anime.AiredOn != nil {
				airedOn := anime.AiredOn.Local()
				animeDto.AiredOn = &airedOn
			}
			if anime.ReleasedOn != nil {
				releasedOn := anime.ReleasedOn.Local()
				animeDto.ReleasedOn = &releasedOn
			}
			if updateErr := sj.AnimeDao.Update(*animeDto); updateErr != nil {
				return nil, errors.Wrap(updateErr, "")
			}
		}
	}
	return animes, nil
}

func rollbackTransaction(tx *sql.Tx, err error) error {
	if rollbackErr := tx.Rollback(); err != nil {
		return errors.Wrap(rollbackErr, "")
	}
	return errors.Wrap(err, "")
}
