package model

const (
	Success   = 0
	ErrorArgs = iota //1
	ErrorValidate

	ErrorUserNotExist
	ErrorRepeatSignUp
	ErrorUserPassWord
	ErrorJournalNotExist
	ErrorLikeAlready
	ErrorLikeNotExist

	ErrorServe = -1
)

var e = map[int]string{
	Success:       "success",
	ErrorServe:    "serve error",
	ErrorArgs:     "error args",
	ErrorValidate: "err validate",

	ErrorUserNotExist:    "user not exist",
	ErrorRepeatSignUp:    "email has sign up",
	ErrorUserPassWord:    "email or password error",
	ErrorJournalNotExist: "journal not exist",
	ErrorLikeAlready:     "already liked",
	ErrorLikeNotExist:    "have not liked",
}

func GetDataMsg(code int) string {
	if str, ok := e[code]; ok {
		return str
	}
	return ""
}
