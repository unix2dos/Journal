package controller

import (
	"Journal/model"

	"log"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {

	args := new(model.SignUpArgs)

	data := NewSetData(c)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		log.Println(err)
		return
	}

	// TODO:验证输入合法

	// 检测是否注册
	log.Println(args.Email)
	//service.MysqlEngine.Where("email=?", args.Email)

}

func Login(c *gin.Context) {

}
