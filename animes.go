package main

import (
	"time"
)

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
