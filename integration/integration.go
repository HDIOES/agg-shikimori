package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ShikimoriJob struct {
	Db *sql.DB
}

func (sj *ShikimoriJob) Run() {
	client := &http.Client{}
	animes := &[]Anime{}
	page := 1
	for len(*animes) == 50 || page == 1 {
		tx, txErr := sj.Db.Begin()
		handleTxError(txErr, tx)
		resp, err := client.Get("https://shikimori.org/api/animes?page=" + strconv.Itoa(page) + "&limit=50")
		handleTxError(err, tx)
		body, err := ioutil.ReadAll(resp.Body)
		handleTxError(err, tx)
		parseError := json.Unmarshal(body, animes)
		handleTxErrorWithAnimesArrays(parseError, tx, animes, &body)
		for i := 0; i < len(*animes); i++ {
			rows, txExecSelectErr := tx.Query("SELECT external_id FROM ANIME WHERE external_id = $1", (*animes)[i].ID)
			handleTxError(txExecSelectErr, tx)
			if !rows.Next() {
				var airedOn *string = nil
				if (*animes)[i].AiredOn != nil {
					airedOn = (*animes)[i].AiredOn.toDateValue()
				}
				var releasedOn *string = nil
				if (*animes)[i].ReleasedOn != nil {
					releasedOn = (*animes)[i].ReleasedOn.toDateValue()
				}
				var posterURL string = (*animes)[i].Image.Original
				_, txExecErr := tx.Exec("INSERT INTO anime (external_id, name, russian, amine_url, kind, anime_status, epizodes, epizodes_aired, aired_on, released_on, poster_url) "+
					"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
					(*animes)[i].ID,
					(*animes)[i].Name,
					(*animes)[i].Russian,
					(*animes)[i].URL,
					(*animes)[i].Kind,
					(*animes)[i].Status,
					(*animes)[i].Episodes,
					(*animes)[i].EpisodesAired,
					airedOn,
					releasedOn,
					posterURL)
				handleTxError(txExecErr, tx)
			}
			rows.Close()
		}
		page++
		handleTxError(tx.Commit(), tx)
		fmt.Println("Page with number " + strconv.Itoa(page) + " has been processed")
		resp.Body.Close()
		time.Sleep(2 * time.Second)
	}
	fmt.Println("Job has been ended")
}

func handleTxError(err error, tx *sql.Tx) {
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Fatal(rollbackErr)
		} else {
			log.Fatal(err)
		}
	}
}

func handleTxErrorWithAnimesArrays(err error, tx *sql.Tx, animes *[]Anime, body *[]byte) {
	if err != nil {
		fmt.Println(string(*body))
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println(rollbackErr)
		} else {
			fmt.Println(err)
		}
	}
}

//Anime struct
type Anime struct {
	ID            int64          `json:"id"`
	Name          string         `json:"name"`
	Russian       string         `json:"russian"`
	Image         Image          `json:"image"`
	URL           string         `json:"url"`
	Kind          string         `json:"kind"`
	Status        string         `json:"status"`
	Episodes      int32          `json:"episodes"`
	EpisodesAired int32          `json:"episodes_aired"`
	AiredOn       *ShikimoriTime `json:"aired_on"`
	ReleasedOn    *ShikimoriTime `json:"released_on"`
}

//Image struct
type Image struct {
	Original string `json:"original"`
	Preview  string `json:"preview"`
	X96      string `json:"x96"`
	X48      string `json:"x48"`
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
