package controller

import (
	"Journal/model"
	"Journal/service"

	"github.com/gin-gonic/gin"
)

func GetInfo(c *gin.Context) {
}

func Signup(c *gin.Context) {

	data := GetData(c)

	args := new(model.SignUpArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		return
	}

	// TODO:验证输入合法

	user := new(model.User)
	// 检测用户是否存在
	exist, _ := service.MysqlEngine.Where("email = ?", args.Email).Get(user)
	if exist {
		data.Ret = model.ErrorRepeatSignUp
		return
	}

	// 存续到数据库
	user.Id = service.GetSnowFlakeId()
	user.Alias = args.Alias
	user.Email = args.Email
	user.Password = args.Password
	userService.SetUserToMysqlAndRedis(user)
	c.Set("uid", user.Id)

	// 存储session
	SessionSave(c)
}

func Login(c *gin.Context) {

	data := GetData(c)

	args := new(model.LoginArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		return
	}

	user := new(model.User)
	//检测用户是否存在
	exist, _ := service.MysqlEngine.Where("email = ?", args.Email).Get(user)
	if !exist {
		data.Ret = model.ErrorUserPassWord
		return
	}

	//检测密码是否正确
	if user.Password != args.Password {
		data.Ret = model.ErrorUserPassWord
		return
	}

	c.Set("uid", user.Id)

	// 存储session
	SessionSave(c)
}
