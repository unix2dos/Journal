package model

const (
	Success   = 0
	ErrorArgs = iota
	ErrorSignUp
)

var e = map[int]string{
	Success:     "success",
	ErrorArgs:   "参数错误",
	ErrorSignUp: "重复注册",
}

func GetDataMsg(code int) string {
	if str, ok := e[code]; ok {
		return str
	}
	return ""
}
