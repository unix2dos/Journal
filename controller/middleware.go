package controller

import (
	"net/http"

	"Journal/model"

	"github.com/gin-gonic/gin"
)

func CommonReturn(c *gin.Context) {
	c.Next() //把hanlder里面的东西执行完了

	data := GetData(c)

	if data.Ret == model.Success {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusInternalServerError, data)
	}

}
