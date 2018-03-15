package model

import "time"

type Journal struct {
	Id        int64  `json:"id,string" xorm:"pk BIGINT(20)"`
	Title     string `json:"title" xorm:"VARCHAR(255)"`
	Content   string `json:"content" xorm:"Text"`
	Public    string `json:"public" xorm:"VARCHAR(20)"`
	LikeCount int64  `json:"like_count,string"`
	//LikeByMe  string `json:"like_by_me" xorm:"-" redis:"-"` //TODO: 考虑太多头疼,回头再考虑这个字段
	UserId int64 `json:"-" xorm:"BIGINT(20)"`

	Create     time.Time `json:"-"  xorm:"DATETIME"`                    //存储数据库的
	Update     time.Time `json:"-"  xorm:"DATETIME"`                    //存储数据库的
	CreateTime int64     `json:"create_time,string" xorm:"-" redis:"-"` //转换给客户端的
	UpdateTime int64     `json:"update_time,string" xorm:"-" redis:"-"` //转换给客户端的
}
