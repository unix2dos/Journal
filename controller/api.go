package controller

import (
	"Journal/model"
	"Journal/service"

	"time"

	"github.com/gin-gonic/gin"
)

func GetInfo(c *gin.Context) {
}

func Signup(c *gin.Context) {

	data := GetData(c)

	args := new(model.SignUpArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("Signup err=%v", err)
		return
	}

	// TODO:验证输入合法

	user := new(model.User)
	// 检测用户是否存在
	exist, _ := service.MysqlEngine.Where("email = ?", args.Email).Get(user)
	if exist {
		service.Logs.Errorf("Signup ErrorRepeatSignUp")
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
		service.Logs.Errorf("Login err=%v", err)
		return
	}

	user := new(model.User)
	//检测用户是否存在
	exist, _ := service.MysqlEngine.Where("email = ?", args.Email).Get(user)
	if !exist {
		service.Logs.Errorf("Login ErrorUserPassWord")
		data.Ret = model.ErrorUserPassWord
		return
	}

	//检测密码是否正确
	if user.Password != args.Password {
		data.Ret = model.ErrorUserPassWord
		service.Logs.Errorf("Login ErrorUserPassWord")
		return
	}

	c.Set("uid", user.Id)

	// 存储session
	SessionSave(c)
}

func JournalList(c *gin.Context) {
	uid, _ := c.Get("uid")
	data := GetData(c)

	list, err := journalService.GetJournalList(uid.(int64)) //TODO: 从redis查?
	if err != nil {
		service.Logs.Errorf("JournalList err=%v", err)
		return
	}

	for _, v := range list {
		v.CreateTime = v.Create.Unix()
		v.UpdateTime = v.Update.Unix()
	}

	data.Data["journals"] = list
}

func JournalAdd(c *gin.Context) {
	data := GetData(c)
	args := new(model.JournalArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("JournalAdd err=%v", err)
		return
	}

	uid, _ := c.Get("uid")
	userId := uid.(int64)

	journal := new(model.Journal)
	journal.Id = service.GetSnowFlakeId()
	journal.Title = args.Title
	journal.Content = args.Content
	journal.Public = args.Public
	journal.LikeCount = 0
	journal.UserId = userId
	journal.Create = time.Now()
	journal.Update = time.Now()
	journal.CreateTime = journal.Create.Unix()
	journal.UpdateTime = journal.Update.Unix()

	if err := journalService.JournalAdd(journal); err != nil {
		service.Logs.Errorf("JournalAdd err=%v", err)
		return
	}

	data.Data["journal"] = journal

}

func JournalUpdate(c *gin.Context) {

}

func JournalDel(c *gin.Context) {

}
