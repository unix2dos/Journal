package service

import (
	"Journal/model"
	"fmt"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) GetUserById(userId int64) (user *model.User, exist bool, err error) {

	key := u.GetUserRedisKey(userId)
	user = new(model.User)

	exist, err = RedisStore.EXISTS(key)
	if err != nil {
		return
	}

	if exist {
		//从redis找
		if err = RedisStore.HMGetStruct(key, user); err != nil {
			return
		}
		return user, true, nil

	} else {
		//从数据库找
		exist, err = MysqlEngine.Id(userId).Get(user)
		if err != nil {
			return
		}

		if exist {
			//写入redis
			u.SetUserToReids(user)
			return user, true, nil
		}
	}

	return
}

func (u *User) SetUserToReids(user *model.User) (err error) {
	key := u.GetUserRedisKey(user.Id)
	return RedisStore.HMSet(key, user)
}

func (u *User) SetUserToMysqlAndRedis(user *model.User) (err error) {
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

	_, err = session.Insert(user)
	if err != nil { //insert错误, update一下, 如果update也错误,彻底错误
		_, err = session.ID(user.Id).Update(user)
		if err != nil {
			return
		}
	}

	err = u.SetUserToReids(user)
	return
}

func (u *User) GetUserRedisKey(id int64) (key string) {
	return fmt.Sprintf(model.RedisKeyUser, id)
}
