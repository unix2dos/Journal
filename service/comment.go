package service

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/unix2dos/Journal/model"
)

type Comment struct {
}

func NewComment() *Comment {
	return &Comment{}
}

//把user_id 换成  user_alias
//把 ReplyCommentId 找到 commit_id, commit_id找到user_id, user_id再换成user_alias
//把user_id 置空
//把ReplyCommentId 置空
func (c *Comment) SetClientAlias(list []*model.Comment) {

	if len(list) <= 0 {
		return
	}
	userMap := make(map[int64]string, 0)
	replyMap := make(map[int64]string, 0)

	//加上用户名
	var userIds string
	for _, v := range list {
		userIds += strconv.Itoa(int(v.UserId)) + ","
	}
	userIds = userIds[:len(userIds)-1]
	aliases, _ := MysqlEngine.SQL("SELECT id, alias FROM user WHERE id in (" + userIds + ")").QueryString()
	for _, v := range aliases {
		id, _ := strconv.ParseInt(v["id"], 10, 64)
		userMap[id] = v["alias"]
	}

	for _, v := range list {
		v.UserAlias = userMap[v.UserId]
		v.UserId = 0
	}

	//加上回复者用户名 TODO: 可以做个优化, 没有ReplyCommentId不走这个
	var replyIds string
	for _, v := range list {
		replyIds += strconv.Itoa(int(v.ReplyCommentId)) + ","
	}
	replyIds = replyIds[:len(replyIds)-1]

	sql := "SELECT a.alias,b.id FROM USER AS a RIGHT JOIN"
	sql += "(SELECT id,user_id FROM COMMENT WHERE id IN ("
	sql += replyIds
	sql += ")) AS b "
	sql += "ON a.id=b.user_id"

	aliases, _ = MysqlEngine.SQL(sql).QueryString()
	for _, v := range aliases {
		id, _ := strconv.ParseInt(v["id"], 10, 64)
		replyMap[id] = v["alias"]
	}
	for _, v := range list {
		v.ReplyUserAlias = replyMap[v.ReplyCommentId]
		v.ReplyCommentId = 0
	}
}

func (c *Comment) GetCommentList(journalId int64) (list []*model.Comment, err error) {
	list = make([]*model.Comment, 0)
	MysqlEngine.Where("journal_id=?", journalId).Find(&list)
	return
}

func (c *Comment) GetCommentById(commentId int64) (comment *model.Comment, exist bool, err error) {

	key := c.getCommentRedisKey(commentId)
	comment = new(model.Comment)

	exist, err = RedisStore.EXISTS(key)
	if err != nil {
		Logs.Errorf("EXISTS commentId=%d err=%v", commentId, err)
		return
	}

	if exist {
		//从redis找
		var str string
		str, err = RedisStore.Get(key)
		if err != nil {
			Logs.Errorf("Get commentId=%d err=%v", commentId, err)
			return
		}
		err = json.Unmarshal([]byte(str), comment)
		if err != nil {
			Logs.Errorf("Unmarshal commentId=%d err=%v", commentId, err)
			return
		}
		return comment, true, nil

	} else {
		//从数据库找
		exist, err = MysqlEngine.Id(commentId).Get(comment)
		if err != nil {
			Logs.Errorf("MysqlEngine Get commentId=%d err=%v", commentId, err)
			return
		}

		if exist {
			//写入redis
			c.SetCommentToReids(comment)
			return comment, true, nil
		}
	}

	return
}

func (c *Comment) SetCommentToReids(comment *model.Comment) (err error) {
	key := c.getCommentRedisKey(comment.Id)
	bytes, err := json.Marshal(comment)
	if err != nil {
		return
	}
	return RedisStore.Set(key, bytes)
}

func (c *Comment) DelCommentFromReids(comment *model.Comment) (err error) {
	key := c.getCommentRedisKey(comment.Id)
	_, err = RedisStore.Del(key)
	return
}

func (c *Comment) SetCommentToMysqlAndRedis(comment *model.Comment) (err error) {
	session := MysqlEngine.NewSession()
	session.Begin()
	defer func() {
		if err == nil {
			session.Commit()
		} else {
			session.Rollback()
		}
		session.Close()
	}()

	if comment.LikeUsers == nil {
		comment.LikeUsers = []int64{}
	}
	_, err = session.Insert(comment)
	if err != nil {
		session = session.MustCols("like_users")
		_, err = session.ID(comment.Id).Update(comment)
		if err != nil {
			return
		}
	}

	err = c.SetCommentToReids(comment)
	return
}

func (c *Comment) DelCommentFromMysqlAndRedis(comment *model.Comment) (err error) {
	session := MysqlEngine.NewSession()
	session.Begin()
	defer func() {
		if err == nil {
			session.Commit()
		} else {
			session.Rollback()
		}
		session.Close()
	}()

	_, err = session.ID(comment.Id).Delete(comment)
	if err != nil {
		return
	}

	err = c.DelCommentFromReids(comment)

	return
}

//--------------------------------------------------//
func (c *Comment) getCommentRedisKey(commentId int64) string {
	return fmt.Sprintf(model.RedisKeyComment, commentId)
}
