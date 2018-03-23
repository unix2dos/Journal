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
	err := UserService.SetUserToMysqlAndRedis(user)
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

	list, err := JournalService.GetJournalList(GetUid(c))
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalList err=%v", err)
		return
	}

	for _, v := range list {
		JournalService.SetClientLikeInfo(GetUid(c), v)
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

	if err := JournalService.SetJournalToMysqlAndRedis(journal); err != nil {
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
	journal, exist, err := JournalService.GetUserJournalById(GetUid(c), id)
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
	if err := JournalService.SetJournalToMysqlAndRedis(journal); err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalUpdate sql err=%v", err)
		return
	}
	JournalService.SetClientLikeInfo(GetUid(c), journal)
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
	journal, exist, err := JournalService.GetUserJournalById(GetUid(c), id)
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

	if err := JournalService.DelJournalFromMysqlAndRedis(journal); err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalDel sql err=%v", err)
		return
	}
}

func JournalRecommend(c *gin.Context) {
	data := GetData(c)
	list, err := JournalService.GetJournalRecommend(GetUid(c))
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalRecommend err=%v", err)
		return
	}
	for _, v := range list {
		JournalService.SetClientLikeInfo(GetUid(c), v)
	}
	data.Data["journals"] = list
}

func CommentList(c *gin.Context) {
	data := GetData(c)
	args := new(model.CommentListArgs)
	if err := c.BindQuery(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("CommentList err=%v", err)
		return
	}

	//先看journal是否存在
	_, exist, err := JournalService.GetJournalById(args.JournalId)
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("CommentList err=%v", err)
		return
	}

	if !exist {
		data.Ret = model.ErrorJournalNotExist
		service.Logs.Errorf("CommentList not exist %v", args.JournalId)
		return
	}

	list, err := CommentService.GetCommentList(args.JournalId)
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("GetCommentList err=%v", err)
		return
	}

	CommentService.SetClientAlias(list)
	data.Data["comments"] = list
}

func CommentAdd(c *gin.Context) {
	data := GetData(c)
	args := new(model.CommentAddArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("CommentAdd err=%v", err)
		return
	}

	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("CommentAdd validate err=%v", err)
		return
	}

	//判断ReplyCommentId是否存在
	if args.ReplyCommentId != 0 {
		_, exist, err := CommentService.GetCommentById(args.ReplyCommentId)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("CommentAdd err=%v", err)
			return
		}

		if !exist {
			data.Ret = model.ErrorCommentNotExist
			service.Logs.Errorf("CommentAdd not exist %v", args.ReplyCommentId)
			return
		}
	}

	comment := new(model.Comment)
	comment.Id = service.GetSnowFlakeId()
	comment.Content = args.Content
	comment.ReplyCommentId = args.ReplyCommentId
	comment.JournalId = args.JournalId
	comment.UserId = GetUid(c)
	comment.CreateTime = model.Time(time.Now())
	comment.UpdateTime = model.Time(time.Now())

	if err := CommentService.SetCommentToMysqlAndRedis(comment); err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("CommentAdd sql err=%v", err)
		return
	}

	CommentService.SetClientAlias([]*model.Comment{comment})
	data.Data["comment"] = comment
}

func CommentUpdate(c *gin.Context) {
	data := GetData(c)
	args := new(model.CommentUpdateArgs)

	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("CommentUpdate err=%v", err)
		return
	}
	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("CommentUpdate validate err=%v", err)
		return
	}

	//判断comment是否存在
	comment, exist, err := CommentService.GetCommentById(args.CommentId)
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("CommentUpdate err=%v", err)
		return
	}

	if !exist {
		data.Ret = model.ErrorCommentNotExist
		service.Logs.Errorf("CommentUpdate not exist %v", args.CommentId)
		return
	}

	comment.Content = args.Content
	comment.UpdateTime = model.Time(time.Now())
	if err := CommentService.SetCommentToMysqlAndRedis(comment); err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("CommentUpdate sql err=%v", err)
		return
	}

	CommentService.SetClientAlias([]*model.Comment{comment})
	data.Data["comment"] = comment
}

func CommentDelete(c *gin.Context) {
	data := GetData(c)
	args := new(model.CommentDeleteArgs)
	if err := c.BindJSON(args); err != nil {
		data.Ret = model.ErrorArgs
		service.Logs.Errorf("CommentDelete err=%v", err)
		return
	}

	if err := service.Validate.Struct(args); err != nil {
		data.Ret = model.ErrorValidate
		service.Logs.Errorf("CommentDelete validate err=%v", err)
		return
	}

	//先看comment是否存在
	comment, exist, err := CommentService.GetCommentById(args.CommentId)
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("CommentDelete err=%v", err)
		return
	}

	if !exist {
		data.Ret = model.ErrorCommentNotExist
		service.Logs.Errorf("CommentDelete not exist %v", args.CommentId)
		return
	}

	if err := CommentService.DelCommentFromMysqlAndRedis(comment); err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("JournalDel sql err=%v", err)
		return
	}
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
		journal, exist, err := JournalService.GetJournalById(args.LikeId)
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
		user, _, _ := UserService.GetUserById(GetUid(c))

		//判断用户是否点赞过
		isJournal := utils.IntContains(journal.LikeUsers, GetUid(c))
		isUser := utils.IntContains(user.LikeJournals, args.LikeId)
		if isJournal && isUser { //这里用&& 防止数据不统一
			data.Ret = model.ErrorLikeAlready
			service.Logs.Errorf("LikeAdd ErrorLikeAlready")
			return
		}

		//添加到数据库
		if !isJournal {
			journal.LikeUsers = append(journal.LikeUsers, GetUid(c))
			err = JournalService.SetJournalToMysqlAndRedis(journal)
			if err != nil {
				data.Ret = model.ErrorServe
				service.Logs.Errorf("LikeAdd SetJournalToMysqlAndRedis err=%v", err)
				return
			}
		}

		//用户也要添加
		if !isUser {
			user.LikeJournals = append(user.LikeJournals, args.LikeId)
			err = UserService.SetUserToMysqlAndRedis(user)
			if err != nil {
				data.Ret = model.ErrorServe
				service.Logs.Errorf("LikeAdd SetUserToMysqlAndRedis err=%v", err)
				return
			}
		}

	} else if args.LikeType == "2" {
		//判断是否有这个comment
		comment, exist, err := CommentService.GetCommentById(args.LikeId)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("LikeAdd GetCommentById err=%v", err)
			return
		}
		if !exist {
			data.Ret = model.ErrorCommentNotExist
			service.Logs.Errorf("LikeAdd ErrorCommentNotExist")
			return
		}
		user, _, _ := UserService.GetUserById(GetUid(c))

		//判断用户是否点赞过
		isComment := utils.IntContains(comment.LikeUsers, GetUid(c))
		isUser := utils.IntContains(user.LikeComments, args.LikeId)
		if isComment && isUser { //这里用&& 防止数据不统一
			data.Ret = model.ErrorLikeAlready
			service.Logs.Errorf("LikeAdd ErrorLikeAlready")
			return
		}

		//添加到数据库
		if !isComment {
			comment.LikeUsers = append(comment.LikeUsers, GetUid(c))
			err = CommentService.SetCommentToMysqlAndRedis(comment)
			if err != nil {
				data.Ret = model.ErrorServe
				service.Logs.Errorf("LikeAdd SetCommentToMysqlAndRedis err=%v", err)
				return
			}
		}

		//用户也要添加
		if !isUser {
			user.LikeComments = append(user.LikeComments, args.LikeId)
			err = UserService.SetUserToMysqlAndRedis(user)
			if err != nil {
				data.Ret = model.ErrorServe
				service.Logs.Errorf("LikeAdd SetUserToMysqlAndRedis err=%v", err)
				return
			}
		}
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
		journal, exist, err := JournalService.GetJournalById(args.LikeId)
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
		user, _, _ := UserService.GetUserById(GetUid(c))

		//判断用户是否点赞过
		isJournal := utils.IntContains(journal.LikeUsers, GetUid(c))
		isUser := utils.IntContains(user.LikeJournals, args.LikeId)
		if !isJournal && !isUser { //这里用&& 防止数据不统一
			data.Ret = model.ErrorLikeNotExist
			service.Logs.Errorf("LikeDelete ErrorLikeNotExist")
			return
		}

		//删除数据库
		if isJournal {
			journal.LikeUsers = utils.SliceRemoveValue(journal.LikeUsers, GetUid(c))
			err = JournalService.SetJournalToMysqlAndRedis(journal)
			if err != nil {
				data.Ret = model.ErrorServe
				service.Logs.Errorf("LikeDelete SetJournalToMysqlAndRedis err=%v", err)
				return
			}
		}

		//用户也要删除
		if isUser {
			user.LikeJournals = utils.SliceRemoveValue(user.LikeJournals, args.LikeId)
			err = UserService.SetUserToMysqlAndRedis(user)
			if err != nil {
				data.Ret = model.ErrorServe
				service.Logs.Errorf("LikeDelete SetUserToMysqlAndRedis err=%v", err)
				return
			}
		}

	} else if args.LikeType == "2" {
		//判断是否有这个comment
		comment, exist, err := CommentService.GetCommentById(args.LikeId)
		if err != nil {
			data.Ret = model.ErrorServe
			service.Logs.Errorf("LikeDelete GetCommentById err=%v", err)
			return
		}
		if !exist {
			data.Ret = model.ErrorCommentNotExist
			service.Logs.Errorf("LikeDelete ErrorCommentNotExist")
			return
		}
		user, _, _ := UserService.GetUserById(GetUid(c))

		//判断用户是否点赞过
		isComment := utils.IntContains(comment.LikeUsers, GetUid(c))
		isUser := utils.IntContains(user.LikeComments, args.LikeId)
		if !isComment && !isUser { //这里用&& 防止数据不统一
			data.Ret = model.ErrorLikeNotExist
			service.Logs.Errorf("LikeDelete ErrorLikeNotExist")
			return
		}

		//删除数据库
		if isComment {
			comment.LikeUsers = utils.SliceRemoveValue(comment.LikeUsers, GetUid(c))
			err = CommentService.SetCommentToMysqlAndRedis(comment)
			if err != nil {
				data.Ret = model.ErrorServe
				service.Logs.Errorf("LikeDelete SetCommentToMysqlAndRedis err=%v", err)
				return
			}
		}

		//用户也要删除
		if isUser {
			user.LikeComments = utils.SliceRemoveValue(user.LikeComments, args.LikeId)
			err = UserService.SetUserToMysqlAndRedis(user)
			if err != nil {
				data.Ret = model.ErrorServe
				service.Logs.Errorf("LikeDelete SetUserToMysqlAndRedis err=%v", err)
				return
			}
		}
	}

}

func ArchiveGet(c *gin.Context) {
	data := GetData(c)
	list, err := JournalService.GetJournalArchive(GetUid(c))
	if err != nil {
		data.Ret = model.ErrorServe
		service.Logs.Errorf("ArchiveGet err=%v", err)
		return
	}
	for _, v := range list {
		JournalService.SetClientLikeInfo(GetUid(c), v)
	}
	data.Data["journals"] = list
}
