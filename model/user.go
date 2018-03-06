package model

type User struct {
	Id       int64  `json:"id,string" xorm:"pk BIGINT(20)"`
	Alias    string `json:"alias" xorm:"VARCHAR(50)"`
	Email    string `json:"email" xorm:"VARCHAR(50)"`
	Password string `json:"-" xorm:"VARCHAR(50)"`
}
