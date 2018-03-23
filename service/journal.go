package service

import (
	"Journal/model"
	"Journal/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

type Journal struct {
}

func NewJournal() *Journal {
	return &Journal{}
}

func (j *Journal) GetJournalList(userId int64) (list []*model.Journal, err error) {
	list = make([]*model.Journal, 0)
	MysqlEngine.Where("user_id=?", userId).Find(&list)
	return
}

func (j *Journal) GetJournalById(journalId int64) (journal *model.Journal, exist bool, err error) {

	key := j.getJournalRedisKey(journalId)
	journal = new(model.Journal)

	exist, err = RedisStore.EXISTS(key)
	if err != nil {
		Logs.Errorf("EXISTS journalId=%d err=%v", journalId, err)
		return
	}

	if exist {
		//从redis找
		var str string
		str, err = RedisStore.Get(key)
		if err != nil {
			Logs.Errorf("Get journalId=%d err=%v", journalId, err)
			return
		}
		err = json.Unmarshal([]byte(str), journal)
		if err != nil {
			Logs.Errorf("Unmarshal journalId=%d err=%v", journalId, err)
			return
		}
		return journal, true, nil

	} else {
		//从数据库找
		exist, err = MysqlEngine.Id(journalId).Get(journal)
		if err != nil {
			Logs.Errorf("MysqlEngine Get journalId=%d err=%v", journalId, err)
			return
		}

		if exist {
			//写入redis
			j.SetJournalToReids(journal)
			return journal, true, nil
		}
	}

	return
}

func (j *Journal) GetUserJournalById(userId int64, journalId int64) (journal *model.Journal, exist bool, err error) {

	journal, exist, err = j.GetJournalById(journalId)
	if exist && journal.UserId != userId {
		exist = false
	}

	return
}

func (j *Journal) SetJournalToReids(journal *model.Journal) (err error) {
	key := j.getJournalRedisKey(journal.Id)
	bytes, err := json.Marshal(journal)
	if err != nil {
		return
	}
	return RedisStore.Set(key, bytes)
}

func (j *Journal) DelJournalFromReids(journal *model.Journal) (err error) {
	key := j.getJournalRedisKey(journal.Id)
	_, err = RedisStore.Del(key)
	return
}

func (j *Journal) SetJournalToMysqlAndRedis(journal *model.Journal) (err error) {
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

	if journal.LikeUsers == nil {
		journal.LikeUsers = []int64{}
	}
	_, err = session.Insert(journal)
	if err != nil {
		session = session.MustCols("like_users")
		_, err = session.ID(journal.Id).Update(journal)
		if err != nil {
			return
		}
	}

	err = j.SetJournalToReids(journal)
	return
}

func (j *Journal) DelJournalFromMysqlAndRedis(journal *model.Journal) (err error) {
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

	_, err = session.ID(journal.Id).Delete(journal)
	if err != nil {
		return
	}

	err = j.DelJournalFromReids(journal)

	return
}

//不是自己的
//不能是公开的
//不是已经喜欢的

//优先点赞数量
//优先最新的
//限制多少条
func (j *Journal) GetJournalRecommend(userId int64) (list []*model.Journal, err error) {
	limit := 10

	list = make([]*model.Journal, 0)
	sql := "SELECT * FROM journal WHERE user_id != ? AND public = ? AND like_users NOT LIKE '%"
	sql += strconv.FormatInt(userId, 10)
	sql += "%'"
	sql += "ORDER BY LENGTH(like_users) DESC, create_time DESC  LIMIT ?"
	MysqlEngine.SQL(sql, userId, "1", limit).Find(&list)

	return
}

func (j *Journal) GetJournalArchive(userId int64) (list []*model.Journal, err error) {
	list = make([]*model.Journal, 0)
	user, _, err := NewUser().GetUserById(userId)
	if err != nil {
		return
	}

	MysqlEngine.In("id", user.LikeJournals).Find(&list)
	return
}

func (j *Journal) SetClientLikeInfo(userId int64, journal *model.Journal) {
	journal.LikeByMe = utils.BoolToString(utils.IntContains(journal.LikeUsers, userId))
	journal.LikeCount = strconv.Itoa(len(journal.LikeUsers))
	journal.LikeUsers = nil
}

//--------------------------------------------------//
func (j *Journal) getJournalRedisKey(journalId int64) string {
	return fmt.Sprintf(model.RedisKeyJournal, journalId)
}
