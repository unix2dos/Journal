package service

import (
	"Journal/model"
	"encoding/json"
	"fmt"
)

type Comment struct {
}

func NewComment() *Comment {
	return &Comment{}
}

func (c *Comment) GetCommentList(journalId int64) (list []*model.Comment, err error) {
	list = make([]*model.Comment, 0) //TODO: 要排序
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
