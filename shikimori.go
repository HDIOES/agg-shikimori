package main

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
	db *sql.DB
}

func (sj *ShikimoriJob) Run() {
	client := &http.Client{}
	animes := &[]Anime{}
	page := 1
	for len(*animes) == 50 || page == 1 {
		tx, txErr := sj.db.Begin()
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
				_, txExecErr := tx.Exec("INSERT INTO anime (external_id, name, russian, amine_url, kind, anime_status, epizodes, epizodes_aired, aired_on, released_on) "+
					"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
					(*animes)[i].ID,
					(*animes)[i].Name,
					(*animes)[i].Russian,
					(*animes)[i].URL,
					(*animes)[i].Kind,
					(*animes)[i].Status,
					(*animes)[i].Episodes,
					(*animes)[i].EpisodesAired,
					airedOn,
					releasedOn)
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
