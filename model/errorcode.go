package model

import "reading-friends-api/model"

const (
	Success   = 0
	ErrorArgs = iota
	ErrorSignUp
	ErrorLogin
)

var e = map[int]string{
	Success:   "success",
	ErrorArgs: "参数错误",
}

func GetErrMsg(data *model.Data, code int) {
	data.Ret = code
	if s, ok := e[code]; ok {
		data.Msg = s
	}
}
