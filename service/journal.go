package service

import (
	"Journal/model"
	"encoding/json"
	"fmt"
)

type Journal struct {
}

func NewJournal() *Journal {
	return &Journal{}
}

func (j *Journal) GetJournalList(userId int64) (list []*model.Journal, err error) {

	list = make([]*model.Journal, 0)
	MysqlEngine.Where("user_id=?", userId).Find(&list)

	//TODO: 是否需要从redis查??
	//通过userId 获取 journals ids
	//根据ids 获取 journal
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

func (j *Journal) SetJournalToReids(journal *model.Journal) (err error) {
	key := j.getJournalRedisKey(journal.Id)
	bytes, err := json.Marshal(journal)
	if err != nil {
		return
	}
	return RedisStore.Set(key, bytes)
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

	_, err = session.Insert(journal)
	if err != nil {
		_, err = session.ID(journal.Id).Update(journal)
		if err != nil {
			return
		}
	}

	err = j.SetJournalToReids(journal)
	return
}

//--------------------------------------------------//
func (j *Journal) getJournalRedisKey(journalId int64) string {
	return fmt.Sprintf(model.RedisKeyJournal, journalId)
}
