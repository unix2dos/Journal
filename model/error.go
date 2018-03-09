package model

const (
	Success    = 0
	ErrorServe = iota
	ErrorArgs
	ErrorNotLogin
	ErrorRepeatSignUp
	ErrorUserPassWord
)

var e = map[int]string{
	Success:           "success",
	ErrorArgs:         "error args",
	ErrorNotLogin:     "not login",
	ErrorRepeatSignUp: "email has sign up",
	ErrorUserPassWord: "email or password error",
}

func GetDataMsg(code int) string {
	if str, ok := e[code]; ok {
		return str
	}
	return ""
}
