package controller

import (
	"Journal/model"

	"github.com/gin-gonic/gin"
)

func GetData(c *gin.Context) *model.Data {
	data, exist := c.Get("data")
	if exist {
		return data.(*model.Data)
	} else {
		return model.NewData()
	}
}

func NewSetData(c *gin.Context) *model.Data {
	data := model.NewData()
	c.Set("data", data)
	return data
}
