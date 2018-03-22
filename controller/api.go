package controller

import (
	"Journal/model"
	"Journal/service"
	"Journal/utils"

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
		service.Logs.Errorf("Signup bind err=%v", err)
		return
	}

	// 验证输入合法
	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("Signup validate err=%v", err)
		return
	}

	// 检测用户是否存在
	user := new(model.User)
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
	user.Password = utils.ScryptPassWord(args.Password)
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

	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("Login validate err=%v", err)
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
	if user.Password != utils.ScryptPassWord(args.Password) {
		data.Ret = model.ErrorUserPassWord
		service.Logs.Errorf("Login ErrorUserPassWord")
		return
	}

	c.Set("uid", user.Id)
	// 存储session
	SessionSave(c)
}

func JournalList(c *gin.Context) {
	data := GetData(c)

	list, err := journalService.GetJournalList(GetUid(c))
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalList err=%v", err)
		return
	}

	for _, v := range list {
		journalService.SetClientLikeInfo(GetUid(c), v)
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

	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("JournalAdd validate err=%v", err)
		return
	}

	journal := new(model.Journal)
	journal.Id = service.GetSnowFlakeId()
	journal.Title = args.Title
	journal.Content = args.Content
	journal.Public = args.Public
	journal.UserId = GetUid(c)
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

	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("JournalUpdate validate err=%v", err)
		return
	}

	//先看journal是否存在
	id, _ := strconv.ParseInt(args.Id, 10, 64)
	journal, exist, err := journalService.GetUserJournalById(GetUid(c), id)
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
	journalService.SetClientLikeInfo(GetUid(c), journal)
	data.Data["journal"] = journal
}

func JournalDel(c *gin.Context) {
	data := GetData(c)
	args := new(model.JournalDeleteArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("JournalDel err=%v", err)
		return
	}

	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("JournalDel validate err=%v", err)
		return
	}

	//先看journal是否存在
	id, _ := strconv.ParseInt(args.Id, 10, 64)
	journal, exist, err := journalService.GetUserJournalById(GetUid(c), id)
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalDel err=%v", err)
		return
	}

	if !exist {
		data.Ret = model.ErrorJournalNotExist
		service.Logs.Errorf("JournalDel not exist %v", args.Id)
		return
	}

	if err := journalService.DelJournalFromMysqlAndRedis(journal); err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalDel sql err=%v", err)
		return
	}
}

func JournalRecommend(c *gin.Context) {
	data := GetData(c)
	list, err := journalService.GetJournalRecommend(GetUid(c))
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalRecommend err=%v", err)
		return
	}
	for _, v := range list {
		journalService.SetClientLikeInfo(GetUid(c), v)
	}
	data.Data["journals"] = list
}

func CommentList(c *gin.Context) {

}

func CommentAdd(c *gin.Context) {

}
func CommentUpdate(c *gin.Context) {

}
func CommentDelete(c *gin.Context) {

}

func LikeAdd(c *gin.Context) {
	data := GetData(c)
	args := new(model.LikeAddArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("LikeAdd err=%v", err)
		return
	}

	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("LikeAdd validate err=%v", err)
		return
	}

	if args.LikeType == "1" {
		//判断是否有这个journal
		journal, exist, err := journalService.GetJournalById(args.LikeId)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("LikeAdd GetJournalById err=%v", err)
			return
		}
		if !exist {
			data.Ret = model.ErrorJournalNotExist
			service.Logs.Errorf("LikeAdd ErrorJournalNotExist")
			return
		}
		user, _, _ := userService.GetUserById(GetUid(c))

		//判断用户是否点赞过 这里用&& 防止数据不统一
		if utils.IntContains(journal.LikeUsers, GetUid(c)) &&
			utils.IntContains(user.LikeJournals, args.LikeId) {
			data.Ret = model.ErrorLikeAlready
			service.Logs.Errorf("LikeAdd ErrorLikeAlready")
			return
		}

		//添加到数据库
		journal.LikeUsers = append(journal.LikeUsers, GetUid(c))
		err = journalService.SetJournalToMysqlAndRedis(journal)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("LikeAdd SetJournalToMysqlAndRedis err=%v", err)
			return
		}
		//用户也要添加
		user.LikeJournals = append(user.LikeJournals, args.LikeId)
		err = userService.SetUserToMysqlAndRedis(user)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("LikeAdd SetUserToMysqlAndRedis err=%v", err)
			return
		}

	} else if args.LikeType == "2" {

	}

}

func LikeDelete(c *gin.Context) {
	data := GetData(c)
	args := new(model.LikeDelArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("LikeDelete err=%v", err)
		return
	}

	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("LikeDelete validate err=%v", err)
		return
	}

	if args.LikeType == "1" {
		//判断是否有这个journal
		journal, exist, err := journalService.GetJournalById(args.LikeId)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("LikeDelete GetJournalById err=%v", err)
			return
		}
		if !exist {
			data.Ret = model.ErrorJournalNotExist
			service.Logs.Errorf("LikeDelete ErrorJournalNotExist")
			return
		}
		user, _, _ := userService.GetUserById(GetUid(c))

		//判断用户是否点赞过, 这里用&& 防止数据不统一
		if !utils.IntContains(journal.LikeUsers, GetUid(c)) &&
			!utils.IntContains(user.LikeJournals, args.LikeId) {
			data.Ret = model.ErrorLikeNotExist
			service.Logs.Errorf("LikeDelete ErrorLikeNotExist")
			return
		}

		//删除数据库
		journal.LikeUsers = utils.SliceRemoveValue(journal.LikeUsers, GetUid(c))
		err = journalService.SetJournalToMysqlAndRedis(journal)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("LikeDelete SetJournalToMysqlAndRedis err=%v", err)
			return
		}
		//用户也要删除
		user.LikeJournals = utils.SliceRemoveValue(user.LikeJournals, args.LikeId)
		err = userService.SetUserToMysqlAndRedis(user)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("LikeDelete SetUserToMysqlAndRedis err=%v", err)
			return
		}

	} else if args.LikeType == "2" {

	}

}
func ArchiveGet(c *gin.Context) {
	data := GetData(c)
	list, err := journalService.GetJournalArchive(GetUid(c))
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("ArchiveGet err=%v", err)
		return
	}
	for _, v := range list {
		journalService.SetClientLikeInfo(GetUid(c), v)
	}
	data.Data["journals"] = list
}
