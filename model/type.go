package model

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(time.Time(*t).Unix(), 10))
}

func (t *Time) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(num, 0))
	return nil
}

func (t Time) String() string {
	return time.Time(t).String()
}
