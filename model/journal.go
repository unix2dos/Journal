package model

import (
	"time"

	"strconv"

	"fmt"

	"github.com/gin-gonic/gin/json"
)

const RedigoTimestamp = "2006-01-02 15:04:05.999999 -0700 MST"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	fmt.Println("MarshalJSON")
	return json.Marshal(strconv.FormatInt(time.Time(t).Unix(), 10))
}

func (t Time) MarshalBinary() ([]byte, error) {
	fmt.Println("MarshalBinary")
	return []byte{'1'}, nil
}

func (t Time) UnmarshalBinary(data []byte) error {
	fmt.Println("UnmarshalBinary")
	return nil
}

func (t Time) MarshalText() ([]byte, error) {
	fmt.Println("MarshalText")
	return []byte{'2'}, nil
}

func (t Time) UnmarshalText(data []byte) error {
	fmt.Println("UnmarshalText")
	return nil
}

type Journal struct {
	Id        int64  `json:"id,string" xorm:"pk BIGINT(20)"`
	Title     string `json:"title" xorm:"VARCHAR(255)"`
	Content   string `json:"content" xorm:"Text"`
	Public    string `json:"public" xorm:"VARCHAR(20)"`
	LikeCount int64  `json:"like_count,string"`
	//LikeByMe  string `json:"like_by_me" xorm:"-" redis:"-"` //TODO: 考虑太多头疼,回头再考虑这个字段
	UserId int64 `json:"-" xorm:"BIGINT(20)"`

	CreateTime Time `json:"create_time"  xorm:"DATETIME"` //1.存数据库, 易看
	UpdateTime Time `json:"update_time"  xorm:"DATETIME"`
}
