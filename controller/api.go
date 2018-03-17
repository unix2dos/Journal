package controller

import (
	"Journal/model"
	"Journal/service"

	"time"

	"strconv"

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
		data.Ret = model.ErrorRepeatSignUp
		service.Logs.Errorf("Signup ErrorRepeatSignUp")
		return
	}

	// 存续到数据库
	user.Id = service.GetSnowFlakeId()
	user.Alias = args.Alias
	user.Email = args.Email
	user.Password = args.Password
	err := userService.SetUserToMysqlAndRedis(user)
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("Signup err=%v", err)
		return
	}

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

	//检测用户是否存在
	user := new(model.User)
	exist, _ := service.MysqlEngine.Where("email = ?", args.Email).Get(user)
	if !exist {
		data.Ret = model.ErrorUserPassWord
		service.Logs.Errorf("Login ErrorUserPassWord")
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
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalList err=%v", err)
		return
	}
	data.Data["journals"] = list
}

func JournalAdd(c *gin.Context) {
	data := GetData(c)
	args := new(model.JournalAddArgs)
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
	journal.CreateTime = model.Time(time.Now())
	journal.UpdateTime = model.Time(time.Now())

	if err := journalService.SetJournalToMysqlAndRedis(journal); err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalAdd sql err=%v", err)
		return
	}

	data.Data["journal"] = journal
}

func JournalUpdate(c *gin.Context) {
	data := GetData(c)
	args := new(model.JournalUpdateArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("JournalUpdate err=%v", err)
		return
	}

	//先看journal是否存在
	id, _ := strconv.ParseInt(args.Id, 10, 64)
	journal, exist, err := journalService.GetJournalById(id)
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalUpdate err=%v", err)
		return
	}

	if !exist {
		data.Ret = model.ErrorJournalNotExist
		service.Logs.Errorf("JournalUpdate not exist %v", args.Id)
		return
	}

	journal.Public = args.Public
	journal.Title = args.Title
	journal.Content = args.Content
	journal.UpdateTime = model.Time(time.Now())
	if err := journalService.SetJournalToMysqlAndRedis(journal); err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalUpdate sql err=%v", err)
		return
	}
	data.Data["journal"] = journal
}

func JournalDel(c *gin.Context) {

}
