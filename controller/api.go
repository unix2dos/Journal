package controller

import (
	"Journal/model"
	"Journal/service"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {

	args := new(model.SignUpArgs)

	data := NewSetData(c)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		return
	}

	// TODO:验证输入合法

	// 检测是否注册
	user := new(model.User)
	//exist, _ := service.MysqlEngine.Where("email = ?", args.Email).Get(user)
	//if exist {
	//	data.Ret = model.ErrorSignUp
	//	return
	//}

	// 存续到数据库
	user.Alias = args.Alias
	user.Email = args.Email
	user.Password = args.Password
	service.MysqlEngine.Insert(user)

}

func Login(c *gin.Context) {

}
