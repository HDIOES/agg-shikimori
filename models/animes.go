package models

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/pkg/errors"
)

//AnimeDAO struct
type AnimeDAO struct {
	Db *sql.DB
}

//DeleteAll function
func (dao *AnimeDAO) DeleteAll() error {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("TRUNCATE anime CASCADE")
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	defer stmt.Close()
	if _, stmtErr := stmt.Exec(); stmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	if cErr := commitTransaction(tx); cErr != nil {
		return errors.Wrap(cErr, "")
	}
	return nil
}

//FindByFilter function
func (dao *AnimeDAO) FindByFilter(filter AnimeQueryBuilder) ([]AnimeDTO, error) {
	query, args := filter.Build()
	stmt, prepareStmtErr := dao.Db.Prepare(query)
	if prepareStmtErr != nil {
		return nil, errors.Wrap(prepareStmtErr, "")
	}
	defer stmt.Close()
	result, stmtErr := stmt.Query(args...)
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, "")
	}
	defer result.Close()
	var capacity int32 = 50
	if filter.Limit > 0 && filter.Limit <= 50 {
		capacity = filter.Limit
	}
	dtos := make([]AnimeDTO, 0, capacity)
	for result.Next() {
		animeDto := AnimeDTO{}
		var id sql.NullInt64
		var name sql.NullString
		var externalID sql.NullString
		var russianName sql.NullString
		var animeURL sql.NullString
		var kind sql.NullString
		var animeStatus sql.NullString
		var epizodes sql.NullInt64
		var epizodesAired sql.NullInt64
		var airedOn sql.NullString
		var releasedOn sql.NullString
		var posterURL sql.NullString
		var score sql.NullFloat64
		var duration sql.NullFloat64
		var rating sql.NullString
		var franchase sql.NullString
		var processed sql.NullBool
		result.Scan(
			&id,
			&name,
			&externalID,
			&russianName,
			&animeURL,
			&kind,
			&animeStatus,
			&epizodes,
			&epizodesAired,
			&airedOn,
			&releasedOn,
			&posterURL,
			&score,
			&duration,
			&rating,
			&franchase,
			&processed)
		if id.Valid {
			animeDto.ID = id.Int64
		}
		if name.Valid {
			animeDto.Name = &name.String
		}
		if externalID.Valid {
			animeDto.ExternalID = externalID.String
		}
		if russianName.Valid {
			animeDto.Russian = &russianName.String
		}
		if animeURL.Valid {
			animeDto.AnimeURL = &animeURL.String
		}
		if kind.Valid {
			animeDto.Kind = &kind.String
		}
		if animeStatus.Valid {
			animeDto.Status = &animeStatus.String
		}
		if epizodes.Valid {
			animeDto.Epizodes = &epizodes.Int64
		}
		if epizodesAired.Valid {
			animeDto.EpizodesAired = &epizodesAired.Int64
		}
		if airedOn.Valid {
			airedOnTime, parseTimeErr := parseTime(airedOn.String)
			if parseTimeErr != nil {
				util.HandleError(parseTimeErr)
			}
			animeDto.AiredOn = &airedOnTime
		}
		if releasedOn.Valid {
			releasedOnTime, parseTimeErr := parseTime(releasedOn.String)
			if parseTimeErr != nil {
				util.HandleError(parseTimeErr)
			}
			animeDto.ReleasedOn = &releasedOnTime
		}
		if posterURL.Valid {
			animeDto.PosterURL = &posterURL.String
		}
		if score.Valid {
			animeDto.Score = &score.Float64
		}
		if duration.Valid {
			animeDto.Duration = &duration.Float64
		}
		if rating.Valid {
			animeDto.Rating = &rating.String
		}
		if franchase.Valid {
			animeDto.Franchise = &franchase.String
		}
		if processed.Valid {
			animeDto.Processed = &processed.Bool
		}
		dtos = append(dtos, animeDto)
	}
	return dtos, nil
}

//FindByExternalID function
func (dao *AnimeDAO) FindByExternalID(externalID string) (*AnimeDTO, error) {
	sqlBuilder := AnimeQueryBuilder{}
	sqlBuilder.AddID(externalID)
	dao.FindByFilter(sqlBuilder)
	animeDtos, findErr := dao.FindByFilter(sqlBuilder)
	if findErr != nil {
		return nil, errors.Wrap(findErr, "")
	}
	if len(animeDtos) > 0 {
		return &animeDtos[0], nil
	}
	return nil, errors.New("Anime not found")
}

//FindByID function
func (dao *AnimeDAO) FindByID(ID int64) (*AnimeDTO, error) {
	stmt, prepareStmtErr := dao.Db.Prepare("SELECT id, name, external_id, russian, amine_url, kind, anime_status, epizodes, epizodes_aired, aired_on, released_on, poster_url, score, duration, rating, franchase FROM anime WHERE id = $1")
	if prepareStmtErr != nil {
		return nil, errors.Wrap(prepareStmtErr, "")
	}
	defer stmt.Close()
	result, stmtErr := stmt.Query(ID)
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, "")
	}
	defer result.Close()
	dto := AnimeDTO{}
	if result.Next() {
		animeDto := AnimeDTO{}
		var id sql.NullInt64
		var name sql.NullString
		var externalID sql.NullString
		var russianName sql.NullString
		var animeURL sql.NullString
		var kind sql.NullString
		var animeStatus sql.NullString
		var epizodes sql.NullInt64
		var epizodesAired sql.NullInt64
		var airedOn sql.NullString
		var releasedOn sql.NullString
		var posterURL sql.NullString
		var score sql.NullFloat64
		var duration sql.NullFloat64
		var rating sql.NullString
		var franchase sql.NullString
		var processed sql.NullBool
		result.Scan(&id,
			&name,
			&externalID,
			&russianName,
			&animeURL,
			&kind,
			&animeStatus,
			&epizodes,
			&epizodesAired,
			&airedOn,
			&releasedOn,
			&posterURL,
			&score,
			&duration,
			&rating,
			&franchase,
			&processed)
		animeDto.ID = id.Int64
		animeDto.Name = &name.String
		animeDto.ExternalID = externalID.String
		animeDto.Russian = &russianName.String
		animeDto.AnimeURL = &animeURL.String
		animeDto.Kind = &kind.String
		animeDto.Status = &animeStatus.String
		animeDto.Epizodes = &epizodes.Int64
		animeDto.EpizodesAired = &epizodesAired.Int64
		airedOnTime, parseTimeErr := parseTime(airedOn.String)
		if parseTimeErr != nil {
			util.HandleError(parseTimeErr)
			animeDto.AiredOn = nil
		} else {
			animeDto.AiredOn = &airedOnTime
		}
		releasedOnTime, parseTimeErr := parseTime(releasedOn.String)
		if parseTimeErr != nil {
			util.HandleError(parseTimeErr)
			animeDto.ReleasedOn = nil
		} else {
			animeDto.ReleasedOn = &releasedOnTime
		}
		animeDto.PosterURL = &posterURL.String
		animeDto.Score = &score.Float64
		animeDto.Duration = &duration.Float64
		animeDto.Rating = &rating.String
		animeDto.Franchise = &franchase.String
		animeDto.Processed = &processed.Bool
	}
	return &dto, nil
}

func parseTime(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", value)
}

//Create finction
func (dao *AnimeDAO) Create(anime AnimeDTO) (int64, error) {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("INSERT INTO anime (external_id, name, russian, amine_url, kind, anime_status, epizodes, epizodes_aired, aired_on, released_on, poster_url, score, duration, rating, franchase, processed, lastmodifytime) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, now()) RETURNING id")
	if prepareStmtErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	defer stmt.Close()
	result, stmtErr := stmt.Query(
		anime.ExternalID,
		anime.Name,
		anime.Russian,
		anime.AnimeURL,
		anime.Kind,
		anime.Status,
		anime.Epizodes,
		anime.EpizodesAired,
		anime.AiredOn,
		anime.ReleasedOn,
		anime.PosterURL,
		anime.Score,
		anime.Duration,
		anime.Rating,
		anime.Franchise,
		anime.Processed)
	if stmtErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	var ID sql.NullInt64
	if result.Next() {
		result.Scan(&ID)
	}
	result.Close()
	if commitErr := tx.Commit(); commitErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return ID.Int64, nil
}

//LinkAnimeAndGenre function
func (dao *AnimeDAO) LinkAnimeAndGenre(animeID int64, genreID int64) error {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("INSERT INTO anime_genre (anime_id, genre_id) VALUES($1, $2)")
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	defer stmt.Close()
	_, execErr := stmt.Exec(animeID, genreID)
	if execErr != nil {
		return rollbackTransaction(tx, errors.Wrap(execErr, ""))
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return nil
}

//LinkAnimeAndStudio function
func (dao *AnimeDAO) LinkAnimeAndStudio(animeID int64, studioID int64) error {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("INSERT INTO anime_studio (anime_id, studio_id) VALUES($1, $2)")
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	defer stmt.Close()
	_, execErr := stmt.Exec(animeID, studioID)
	if execErr != nil {
		return rollbackTransaction(tx, errors.Wrap(execErr, ""))
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return nil
}

//Update function
func (dao *AnimeDAO) Update(anime AnimeDTO) error {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("UPDATE anime ")
	queryBuilder.WriteString("SET external_id = $1, ")
	queryBuilder.WriteString("name = $2, ")
	queryBuilder.WriteString("russian = $3, ")
	queryBuilder.WriteString("amine_url = $4, ")
	queryBuilder.WriteString("kind = $5, ")
	queryBuilder.WriteString("anime_status = $6, ")
	queryBuilder.WriteString("epizodes = $7, ")
	queryBuilder.WriteString("epizodes_aired = $8, ")
	queryBuilder.WriteString("aired_on = $9, ")
	queryBuilder.WriteString("released_on = $10, ")
	queryBuilder.WriteString("poster_url = $11, ")
	queryBuilder.WriteString("score = $12, ")
	queryBuilder.WriteString("duration = $13, ")
	queryBuilder.WriteString("rating = $14, ")
	queryBuilder.WriteString("franchase = $15, ")
	queryBuilder.WriteString("processed = $16, ")
	queryBuilder.WriteString("lastmodifytime = now() ")
	queryBuilder.WriteString("WHERE id = $17")
	stmt, prepareStmtErr := tx.Prepare(queryBuilder.String())
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	_, stmtErr := stmt.Exec(
		anime.ExternalID,
		anime.Name,
		anime.Russian,
		anime.AnimeURL,
		anime.Kind,
		anime.Status,
		anime.Epizodes,
		anime.EpizodesAired,
		anime.AiredOn,
		anime.ReleasedOn,
		anime.PosterURL,
		anime.Score,
		anime.Duration,
		anime.Rating,
		anime.Franchise,
		anime.Processed,
		anime.ID)
	if stmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return nil
}

//GetCount function
func (dao *AnimeDAO) GetCount(sqlBuilder AnimeQueryBuilder) (int64, error) {
	sqlQuery, args := sqlBuilder.Build()
	stmt, prepareStmtErr := dao.Db.Prepare(sqlQuery)
	if prepareStmtErr != nil {
		return 0, errors.Wrap(prepareStmtErr, "")
	}
	defer stmt.Close()
	result, queryErr := stmt.Query(args...)
	if queryErr != nil {
		return 0, errors.Wrap(queryErr, "")
	}
	defer result.Close()
	if result.Next() {
		var count sql.NullInt64
		result.Scan(&count)
		return count.Int64, nil
	}
	return 0, nil
}

//GetRandomAnime function returns random anime by predefined filter
func (dao *AnimeDAO) GetRandomAnime(sqlBuilder AnimeQueryBuilder) (*AnimeDTO, error) {
	query, args := sqlBuilder.Build()
	stmt, prepareStmtErr := dao.Db.Prepare(query)
	if prepareStmtErr != nil {
		return nil, errors.Wrap(prepareStmtErr, "")
	}
	defer stmt.Close()
	result, stmtErr := stmt.Query(args...)
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, "")
	}
	defer result.Close()
	animeDto := AnimeDTO{}
	if result.Next() {
		var id sql.NullInt64
		var name sql.NullString
		var externalID sql.NullString
		var russianName sql.NullString
		var animeURL sql.NullString
		var kind sql.NullString
		var animeStatus sql.NullString
		var epizodes sql.NullInt64
		var epizodesAired sql.NullInt64
		var airedOn sql.NullString
		var releasedOn sql.NullString
		var posterURL sql.NullString
		var score sql.NullFloat64
		var duration sql.NullFloat64
		var rating sql.NullString
		var franchase sql.NullString
		var processed sql.NullBool
		result.Scan(&id,
			&name,
			&externalID,
			&russianName,
			&animeURL,
			&kind,
			&animeStatus,
			&epizodes,
			&epizodesAired,
			&airedOn,
			&releasedOn,
			&posterURL,
			&score,
			&duration,
			&rating,
			&franchase,
			&processed)
		if id.Valid {
			animeDto.ID = id.Int64
		}
		if name.Valid {
			animeDto.Name = &name.String
		}
		if externalID.Valid {
			animeDto.ExternalID = externalID.String
		}
		if russianName.Valid {
			animeDto.Russian = &russianName.String
		}
		if animeURL.Valid {
			animeDto.AnimeURL = &animeURL.String
		}
		if kind.Valid {
			animeDto.Kind = &kind.String
		}
		if animeStatus.Valid {
			animeDto.Status = &animeStatus.String
		}
		if epizodes.Valid {
			animeDto.Epizodes = &epizodes.Int64
		}
		if epizodesAired.Valid {
			animeDto.EpizodesAired = &epizodesAired.Int64
		}
		if airedOn.Valid {
			airedOnTime, parseTimeErr := parseTime(airedOn.String)
			if parseTimeErr != nil {
				util.HandleError(parseTimeErr)
			}
			animeDto.AiredOn = &airedOnTime
		}
		if releasedOn.Valid {
			releasedOnTime, parseTimeErr := parseTime(releasedOn.String)
			if parseTimeErr != nil {
				util.HandleError(parseTimeErr)
			}
			animeDto.ReleasedOn = &releasedOnTime
		}
		if posterURL.Valid {
			animeDto.PosterURL = &posterURL.String
		}
		if score.Valid {
			animeDto.Score = &score.Float64
		}
		if duration.Valid {
			animeDto.Duration = &duration.Float64
		}
		if rating.Valid {
			animeDto.Rating = &rating.String
		}
		if franchase.Valid {
			animeDto.Franchise = &franchase.String
		}
		if processed.Valid {
			animeDto.Processed = &processed.Bool
		}
	}
	return &animeDto, nil
}

//AnimeQueryBuilder struct
type AnimeQueryBuilder struct {
	Limit      int32
	Offset     int32
	Status     string
	Kind       string
	Phrase     string
	Order      string
	Score      int32
	GenreIds   []string
	StudioIds  []string
	Duration   string
	Rating     string
	Franchise  string
	Ids        []string
	ExcludeIds []string
	SQLQuery   strings.Builder
	Processed  bool
	CountOnly  bool
	RowNumber  int64
}

//SetProcessed function
func (aqb *AnimeQueryBuilder) SetProcessed(processed bool) {
	aqb.Processed = processed
}

//AddExcludeID function
func (aqb *AnimeQueryBuilder) AddExcludeID(excludeID string) {
	aqb.ExcludeIds = append(aqb.ExcludeIds, excludeID)
}

//AddID function
func (aqb *AnimeQueryBuilder) AddID(id string) {
	aqb.Ids = append(aqb.Ids, id)
}

//SetFranchise function
func (aqb *AnimeQueryBuilder) SetFranchise(franchise string) {
	aqb.Franchise = franchise
}

//SetRating function
func (aqb *AnimeQueryBuilder) SetRating(rating string) {
	aqb.Rating = rating
}

//SetDuration function
func (aqb *AnimeQueryBuilder) SetDuration(duration string) {
	aqb.Duration = duration
}

//AddStudioID function
func (aqb *AnimeQueryBuilder) AddStudioID(studioID string) {
	aqb.StudioIds = append(aqb.StudioIds, studioID)
}

//AddGenreID function
func (aqb *AnimeQueryBuilder) AddGenreID(genreID string) {
	aqb.GenreIds = append(aqb.GenreIds, genreID)
}

//SetScore function
func (aqb *AnimeQueryBuilder) SetScore(score int32) {
	aqb.Score = score
}

//SetOrder function
func (aqb *AnimeQueryBuilder) SetOrder(order string) {
	aqb.Order = order
}

//SetPhrase function
func (aqb *AnimeQueryBuilder) SetPhrase(phrase string) {
	aqb.Phrase = phrase
}

//SetKind function
func (aqb *AnimeQueryBuilder) SetKind(kind string) {
	aqb.Kind = kind
}

//SetStatus function
func (aqb *AnimeQueryBuilder) SetStatus(status string) {
	aqb.Status = status
}

//SetLimit func
func (aqb *AnimeQueryBuilder) SetLimit(limit int32) {
	aqb.Limit = limit
}

//SetOffset func
func (aqb *AnimeQueryBuilder) SetOffset(offset int32) {
	aqb.Offset = offset
}

//SetCountOnly function
func (aqb *AnimeQueryBuilder) SetCountOnly(countOnly bool) {
	aqb.CountOnly = countOnly
}

//SetRowNumber function
func (aqb *AnimeQueryBuilder) SetRowNumber(rowNumber int64) {
	aqb.RowNumber = rowNumber
}

//Build function
func (aqb *AnimeQueryBuilder) Build() (string, []interface{}) {
	if aqb.CountOnly {
		aqb.SQLQuery.WriteString("SELECT COUNT(DISTINCT ")
		aqb.SQLQuery.WriteString("query.anime_internal_id,")
		aqb.SQLQuery.WriteString("query.name,")
		aqb.SQLQuery.WriteString("query.anime_external_id,")
		aqb.SQLQuery.WriteString("query.russian,")
		aqb.SQLQuery.WriteString("query.amine_url,")
		aqb.SQLQuery.WriteString("query.kind,")
		aqb.SQLQuery.WriteString("query.anime_status,")
		aqb.SQLQuery.WriteString("query.epizodes,")
		aqb.SQLQuery.WriteString("query.epizodes_aired,")
		aqb.SQLQuery.WriteString("query.aired_on,")
		aqb.SQLQuery.WriteString("query.released_on,")
		aqb.SQLQuery.WriteString("query.poster_url,")
		aqb.SQLQuery.WriteString("query.score,")
		aqb.SQLQuery.WriteString("query.duration,")
		aqb.SQLQuery.WriteString("query.rating,")
		aqb.SQLQuery.WriteString("query.franchase,")
		aqb.SQLQuery.WriteString("query.processed) FROM (")
	} else {
		aqb.SQLQuery.WriteString("SELECT DISTINCT ")
		aqb.SQLQuery.WriteString("query.anime_internal_id,")
		aqb.SQLQuery.WriteString("query.name,")
		aqb.SQLQuery.WriteString("query.anime_external_id,")
		aqb.SQLQuery.WriteString("query.russian,")
		aqb.SQLQuery.WriteString("query.amine_url,")
		aqb.SQLQuery.WriteString("query.kind,")
		aqb.SQLQuery.WriteString("query.anime_status,")
		aqb.SQLQuery.WriteString("query.epizodes,")
		aqb.SQLQuery.WriteString("query.epizodes_aired,")
		aqb.SQLQuery.WriteString("query.aired_on,")
		aqb.SQLQuery.WriteString("query.released_on,")
		aqb.SQLQuery.WriteString("query.poster_url,")
		aqb.SQLQuery.WriteString("query.score,")
		aqb.SQLQuery.WriteString("query.duration,")
		aqb.SQLQuery.WriteString("query.rating,")
		aqb.SQLQuery.WriteString("query.franchase,")
		aqb.SQLQuery.WriteString("query.processed ")
		aqb.SQLQuery.WriteString("FROM (")
	}
	aqb.SQLQuery.WriteString("SELECT ")
	aqb.SQLQuery.WriteString("row_number() over(),")
	aqb.SQLQuery.WriteString("animes.anime_internal_id,")
	aqb.SQLQuery.WriteString("animes.name,")
	aqb.SQLQuery.WriteString("animes.anime_external_id,")
	aqb.SQLQuery.WriteString("animes.russian,")
	aqb.SQLQuery.WriteString("animes.amine_url,")
	aqb.SQLQuery.WriteString("animes.kind,")
	aqb.SQLQuery.WriteString("animes.anime_status,")
	aqb.SQLQuery.WriteString("animes.epizodes,")
	aqb.SQLQuery.WriteString("animes.epizodes_aired,")
	aqb.SQLQuery.WriteString("animes.aired_on,")
	aqb.SQLQuery.WriteString("animes.released_on,")
	aqb.SQLQuery.WriteString("animes.poster_url,")
	aqb.SQLQuery.WriteString("animes.score,")
	aqb.SQLQuery.WriteString("animes.duration,")
	aqb.SQLQuery.WriteString("animes.rating,")
	aqb.SQLQuery.WriteString("animes.franchase,")
	aqb.SQLQuery.WriteString("animes.processed ")
	aqb.SQLQuery.WriteString("FROM ")
	aqb.SQLQuery.WriteString("(")
	aqb.SQLQuery.WriteString("SELECT ")
	aqb.SQLQuery.WriteString("anime.id AS anime_internal_id,")
	aqb.SQLQuery.WriteString("anime.name,")
	aqb.SQLQuery.WriteString("anime.external_id as anime_external_id,")
	aqb.SQLQuery.WriteString("anime.russian,")
	aqb.SQLQuery.WriteString("anime.amine_url,")
	aqb.SQLQuery.WriteString("anime.kind,")
	aqb.SQLQuery.WriteString("anime.anime_status,")
	aqb.SQLQuery.WriteString("anime.epizodes,")
	aqb.SQLQuery.WriteString("anime.epizodes_aired,")
	aqb.SQLQuery.WriteString("anime.aired_on,")
	aqb.SQLQuery.WriteString("anime.released_on,")
	aqb.SQLQuery.WriteString("anime.poster_url,")
	aqb.SQLQuery.WriteString("anime.score,")
	aqb.SQLQuery.WriteString("anime.duration,")
	aqb.SQLQuery.WriteString("anime.rating,")
	aqb.SQLQuery.WriteString("anime.franchase,")
	aqb.SQLQuery.WriteString("anime.processed ")
	countOfParameter := 0
	args := make([]interface{}, 0)
	if len(aqb.GenreIds) > 0 {
		aqb.SQLQuery.WriteString(", genre.external_id as genre_external_id")
	}
	if len(aqb.StudioIds) > 0 {
		aqb.SQLQuery.WriteString(", studio.external_id as studio_external_id")
	}
	if len(aqb.Phrase) > 0 {
		countOfParameter++
		aqb.SQLQuery.WriteString(", to_tsvector(anime.russian) as russian_tsvector, to_tsvector(anime.name) as english_tsvector, phraseto_tsquery($")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		aqb.SQLQuery.WriteString(") as ts_query")
		args = append(args, aqb.Phrase)
	}
	aqb.SQLQuery.WriteString(" FROM anime")
	if len(aqb.GenreIds) > 0 {
		aqb.SQLQuery.WriteString(" JOIN anime_genre ON anime.id = anime_genre.anime_id")
		aqb.SQLQuery.WriteString(" JOIN genre ON genre.id = anime_genre.genre_id")
	}
	if len(aqb.StudioIds) > 0 {
		aqb.SQLQuery.WriteString(" JOIN anime_studio ON anime.id = anime_studio.anime_id")
		aqb.SQLQuery.WriteString(" JOIN studio ON studio.id = anime_studio.studio_id")
	}
	aqb.SQLQuery.WriteString(") as animes")
	aqb.SQLQuery.WriteString(" WHERE 1=1")
	if len(aqb.Phrase) > 0 {
		aqb.SQLQuery.WriteString(" AND (animes.russian_tsvector @@ animes.ts_query OR animes.english_tsvector @@ animes.ts_query)")
	}
	if len(aqb.GenreIds) > 0 {
		aqb.SQLQuery.WriteString(" AND genre_external_id IN (")
		for ind, genreExternalID := range aqb.GenreIds {
			countOfParameter++
			args = append(args, genreExternalID)
			aqb.SQLQuery.WriteString("$")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			if ind < len(aqb.GenreIds)-1 {
				aqb.SQLQuery.WriteString(",")
			}
		}
		aqb.SQLQuery.WriteString(")")
	}
	if len(aqb.StudioIds) > 0 {
		aqb.SQLQuery.WriteString(" AND studio_external_id IN (")
		for ind, studioExternalID := range aqb.StudioIds {
			countOfParameter++
			args = append(args, studioExternalID)
			aqb.SQLQuery.WriteString("$")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			if ind < len(aqb.StudioIds)-1 {
				aqb.SQLQuery.WriteString(",")
			}
		}
		aqb.SQLQuery.WriteString(")")
	}
	if len(aqb.Status) > 0 {
		countOfParameter++
		aqb.SQLQuery.WriteString(" AND animes.anime_status = $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.Status)
	}
	if len(aqb.Kind) > 0 {
		var kinds = [...]string{"tv", "movie", "ova", "ona", "special", "music", "tv_13", "tv_24", "tv_48"}
		for _, s := range kinds {
			if strings.Compare(s, aqb.Kind) == 0 {
				countOfParameter++
				aqb.SQLQuery.WriteString(" AND animes.kind = $")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, aqb.Kind)
				break
			}
		}
	}
	if len(aqb.Ids) > 0 {
		aqb.SQLQuery.WriteString(" AND anime_external_id IN (")
		for ind, id := range aqb.Ids {
			countOfParameter++
			args = append(args, id)
			aqb.SQLQuery.WriteString("$")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			if ind < len(aqb.Ids)-1 {
				aqb.SQLQuery.WriteString(",")
			}
		}
		aqb.SQLQuery.WriteString(")")
	}
	if len(aqb.ExcludeIds) > 0 {
		aqb.SQLQuery.WriteString(" AND anime_external_id NOT IN (")
		for ind, excludeID := range aqb.ExcludeIds {
			countOfParameter++
			args = append(args, excludeID)
			aqb.SQLQuery.WriteString("$")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			if ind < len(aqb.ExcludeIds)-1 {
				aqb.SQLQuery.WriteString(",")
			}
		}
		aqb.SQLQuery.WriteString(")")
	}
	if len(aqb.Duration) > 0 {
		switch aqb.Duration {
		case "S":
			{
				aqb.SQLQuery.WriteString(" AND animes.duration <= 10")
			}
		case "D":
			{
				aqb.SQLQuery.WriteString(" AND animes.duration < 30 AND animes.duration >= 10")
			}
		case "F":
			{
				aqb.SQLQuery.WriteString(" AND animes.duration >= 30")
			}
		}
	}
	if len(aqb.Franchise) > 0 {
		countOfParameter++
		aqb.SQLQuery.WriteString(" AND animes.franchase = $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.Franchise)
	}
	if len(aqb.Rating) > 0 {
		var ratings = [...]string{"none", "g", "pg", "pg_13", "r", "r_plus", "rx"}
		for _, r := range ratings {
			if strings.Compare(r, aqb.Rating) == 0 {
				countOfParameter++
				aqb.SQLQuery.WriteString(" AND animes.rating = $")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, aqb.Rating)
				break
			}
		}
	}
	if aqb.Score > 0.0 {
		//need to validate score
		countOfParameter++
		aqb.SQLQuery.WriteString(" AND animes.score >= $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.Score)
	}
	if aqb.Processed == true {
		countOfParameter++
		aqb.SQLQuery.WriteString(" AND animes.processed = $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.Processed)
	}
	if len(aqb.Order) > 0 {
		aqb.SQLQuery.WriteString(" ORDER BY ")
		switch aqb.Order {
		case "id":
			{
				aqb.SQLQuery.WriteString("anime_external_id")
			}
		case "kind":
			{
				aqb.SQLQuery.WriteString("animes.kind")
			}
		case "name":
			{
				aqb.SQLQuery.WriteString("animes.name")
			}
		case "aired_on":
			{
				aqb.SQLQuery.WriteString("animes.aired_on")
			}
		case "episodes":
			{
				aqb.SQLQuery.WriteString("animes.epizodes")
			}
		case "status":
			{
				aqb.SQLQuery.WriteString("animes.status")
			}
		case "relevance":
			{
				if len(aqb.Phrase) > 0 {
					aqb.SQLQuery.WriteString("get_rank(animes.russian_tsvector, animes.english_tsvector, animes.ts_query) DESC")
				}
			}
		}
	}
	if aqb.CountOnly {
		aqb.SQLQuery.WriteString(") as query")
	} else if aqb.RowNumber > 0 {
		countOfParameter++
		aqb.SQLQuery.WriteString(") as query where query.row_number = $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.RowNumber)
	} else {
		if aqb.Limit > 0 {
			countOfParameter++
			aqb.SQLQuery.WriteString(") as query LIMIT $" + strconv.Itoa(countOfParameter))
			args = append(args, aqb.Limit)
		} else {
			aqb.SQLQuery.WriteString(") as query LIMIT 50")
		}
		if aqb.Offset > 0 {
			countOfParameter++
			aqb.SQLQuery.WriteString(" OFFSET $")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			args = append(args, aqb.Offset)
		}
	}
	return aqb.SQLQuery.String(), args
}

//AnimeDTO struct
type AnimeDTO struct {
	ID             int64
	ExternalID     string
	Name           *string
	Russian        *string
	AnimeURL       *string
	Kind           *string
	Status         *string
	Epizodes       *int64
	EpizodesAired  *int64
	AiredOn        *time.Time
	ReleasedOn     *time.Time
	PosterURL      *string
	Score          *float64
	Duration       *float64
	Rating         *string
	Franchise      *string
	Processed      *bool
	LastModifyTime *time.Time
}
