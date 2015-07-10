package model

import (
	"encoding/json"
	"errors"
	"time"
)

const TsLayout = "2006-01-02T15:04:05.000Z"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UTC().Format(TsLayout))
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if data[0] != []byte(`"`)[0] || data[len(data)-1] != []byte(`"`)[0] {
		return errors.New("Not quoted")
	}
	*t, err = ParseTime(string(data[1 : len(data)-1]))
	return
}

func ParseTime(timespec string) (Time, error) {
	t, err := time.Parse(TsLayout, timespec)
	return Time(t), err
}

func MustParseTime(timespec string) Time {
	ts, err := ParseTime(timespec)
	if err != nil {
		panic(err)
	}
	return ts
}
