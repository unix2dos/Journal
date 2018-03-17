package model

const (
	Success   = 0
	ErrorArgs = iota //1
	ErrorNotLogin
	ErrorRepeatSignUp
	ErrorUserPassWord
	ErrorJournalNotExist

	ErrorServe = -1
)

var e = map[int]string{
	Success:              "success",
	ErrorServe:           "serve error",
	ErrorArgs:            "error args",
	ErrorNotLogin:        "not login",
	ErrorRepeatSignUp:    "email has sign up",
	ErrorUserPassWord:    "email or password error",
	ErrorJournalNotExist: "journal not exist",
}

func GetDataMsg(code int) string {
	if str, ok := e[code]; ok {
		return str
	}
	return ""
}
