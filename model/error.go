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
	ErrorArgs:         "参数错误",
	ErrorNotLogin:     "没有登录",
	ErrorRepeatSignUp: "重复注册",
	ErrorUserPassWord: "用户名或密码不正确",
}

func GetDataMsg(code int) string {
	if str, ok := e[code]; ok {
		return str
	}
	return ""
}
