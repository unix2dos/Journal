package controller

import (
	"Journal/model"
	"Journal/service"
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetInfo(c *gin.Context) {
}

func Signup(c *gin.Context) {

	data := GetData(c)

	args := new(model.SignUpArgs)
	if err := c.BindJSON(args); err != nil {
		log.Println("err:", err)
		data.Ret = model.ErrorArgs
		return
	}

	// TODO:验证输入合法

	// 检测是否注册
	user := new(model.User)
	exist, _ := service.MysqlEngine.Where("email = ?", args.Email).Get(user)
	if exist {
		data.Ret = model.ErrorRepeatSignUp
		return
	}

	// 存续到数据库 redis
	user.Id = service.GetSnowFlakeId()
	user.Alias = args.Alias
	user.Email = args.Email
	user.Password = args.Password
	UserSaveDB(user)

	// 存储session
	c.Set("uid", user.Id)
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
	exist, _ := service.MysqlEngine.Where("email = ?", args.Email).Get(user)
	if !exist {
		data.Ret = model.ErrorUserPassWord
		return
	}

	if user.Password != args.Password {
		data.Ret = model.ErrorUserPassWord
		return
	}

	// 存储session
	c.Set("uid", user.Id)
	SessionSave(c)
}

func SessionSave(c *gin.Context) {
	useId, _ := c.Get("uid")
	session := sessions.Default(c)
	session.Set("uid", useId)
	session.Save()
}

func UserSaveDB(user *model.User) (err error) {
	session := service.MysqlEngine.NewSession()
	session.Begin()
	defer func() {
		if err == nil {
			session.Commit()
		} else {
			session.Rollback()
		}
		session.Close()
	}()

	_, err = session.Insert(user)
	if err != nil {
		return
	}

	key := fmt.Sprintf(model.RedisKeyUser, user.Id)
	err = service.RedisStore.HMSet(key, user)
	return
}
