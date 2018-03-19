package model

import (
	"encoding/json"
	"time"

	"strconv"

	"strings"
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

//redis存序列化字符串
type Journal struct {
	Id      int64  `json:"id,string" xorm:"pk BIGINT(20)"`
	Title   string `json:"title" xorm:"VARCHAR(255)"`
	Content string `json:"content" xorm:"Text"`
	Public  string `json:"public" xorm:"VARCHAR(20)"`
	UserId  int64  `json:"-" xorm:"BIGINT(20)"`

	CreateTime Time `json:"create_time"  xorm:"DATETIME"`
	UpdateTime Time `json:"update_time"  xorm:"DATETIME"`

	LikeCount int64  `json:"like_count,string,omitempty" xorm:"-"` //这两个需要自己算
	LikeByMe  string `json:"like_by_me,omitempty" xorm:"-"`
}
