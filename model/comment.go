package model

type Comment struct {
	Id             int64  `json:"id,string" xorm:"pk BIGINT(20)"`
	Content        string `json:"content" xorm:"Text"`
	ReplyCommentId int64  `json:"reply_comment_id,string" xorm:"BIGINT(20) default(0)"`
	UserId         int64  `json:"user_id,string" xorm:"BIGINT(20)"`    //评论作者id
	JournalId      int64  `json:"journal_id,string" xorm:"BIGINT(20)"` //属于哪个journal

	CreateTime Time `json:"create_time"  xorm:"DATETIME"`
	UpdateTime Time `json:"update_time"  xorm:"DATETIME"`
}
