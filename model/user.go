package model

import "time"

type User struct {
	Id       int64     `json:"id,string" xorm:"pk BIGINT(20)"`
	Alias    string    `json:"alias" xorm:"VARCHAR(50)"`
	Email    string    `json:"email" xorm:"VARCHAR(50)"`
	Password string    `json:"-" xorm:"VARCHAR(50)"`
	Create   time.Time `json:"-" redis:"-" xorm:"created DATETIME"`
	Update   time.Time `json:"-" redis:"-" xorm:"updated DATETIME"`
}

type UserLike struct {
	UserId    int64 `json:"user_id,string"`
	JournalId int64 `json:"journal_id,string,omitempty" xorm:"default(0)"`
	CommentId int64 `json:"comment_id,string,omitempty" xorm:"default(0)"`
}
