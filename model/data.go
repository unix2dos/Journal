package model

import "time"

type Data struct {
	Ret  int                    `json:"ret, string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func NewData() *Data {
	d := &Data{}
	d.Ret = Success
	d.Msg = ""
	d.Data = make(map[string]interface{})
	return d
}

type User struct {
	Id       int64  `json:"id,string" xorm:"pk BIGINT(20)"`
	Alias    string `json:"alias" xorm:"VARCHAR(50)"`
	Email    string `json:"email" xorm:"VARCHAR(50)"`
	Password string `json:"-" xorm:"VARCHAR(50)"`

	LikeJournals []int64 `json:"like_journals,omitempty"`
	LikeComments []int64 `json:"like_comments,omitempty"`

	Create time.Time `json:"-" redis:"-" xorm:"created DATETIME"`
	Update time.Time `json:"-" redis:"-" xorm:"updated DATETIME"`
}

//redis存json字符串
type Journal struct {
	Id         int64   `json:"id,string" xorm:"pk BIGINT(20)"`
	Title      string  `json:"title" xorm:"VARCHAR(255)"`
	Content    string  `json:"content" xorm:"Text"`
	Public     string  `json:"public" xorm:"VARCHAR(20)"`
	UserId     int64   `json:"user_id" xorm:"BIGINT(20)"`
	LikeUsers  []int64 `json:"like_users,omitempty"`
	CreateTime Time    `json:"create_time"  xorm:"DATETIME"`
	UpdateTime Time    `json:"update_time"  xorm:"DATETIME"`

	LikeCount int64  `json:"like_count,string" xorm:"-"` //这两个服务器算给客户端
	LikeByMe  string `json:"like_by_me" xorm:"-"`
}

type Comment struct {
	Id             int64   `json:"id,string" xorm:"pk BIGINT(20)"`
	Content        string  `json:"content" xorm:"Text"`
	ReplyCommentId int64   `json:"reply_comment_id,string" xorm:"BIGINT(20) default(0)"`
	UserId         int64   `json:"user_id,string" xorm:"BIGINT(20)"`    //评论作者id
	JournalId      int64   `json:"journal_id,string" xorm:"BIGINT(20)"` //属于哪个journal
	LikeUsers      []int64 `json:"like_users,omitempty"`
	CreateTime     Time    `json:"create_time"  xorm:"DATETIME"`
	UpdateTime     Time    `json:"update_time"  xorm:"DATETIME"`
}
