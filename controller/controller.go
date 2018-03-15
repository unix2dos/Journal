package controller

import (
	"Journal/model"

	"Journal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	userService = service.NewUser()

	NotSessionFilter = []string{
		"/signup",
		"/login",
	}
)

func GetData(c *gin.Context) *model.Data {
	data, exist := c.Get("data")
	if exist {
		return data.(*model.Data)
	} else {
		return SetNewData(c)
	}
}

func SetNewData(c *gin.Context) *model.Data {
	data := model.NewData()
	c.Set("data", data)
	return data
}

func SessionGet(c *gin.Context) (userId int64, ok bool) {
	session := sessions.Default(c)
	userId, ok = session.Get("uid").(int64)
	return
}

func SessionSave(c *gin.Context) {
	useId, _ := c.Get("uid")
	session := sessions.Default(c)
	session.Set("uid", useId)
	session.Save()
}
