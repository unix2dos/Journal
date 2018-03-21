package service

import (
	"Journal/model"
	"io/ioutil"

	"Journal/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"github.com/zheng-ji/goSnowFlake"
	"gopkg.in/go-playground/validator.v8"
)

var (
	MysqlEngine *xorm.Engine
	RedisStore  *utils.RedisStore
	SnowFlake   *goSnowFlake.IdWorker
	Logs        *logrus.Logger
	Validate    *validator.Validate
)

func SInit() {
	LogInit()
	SqlInit()

	var err error

	//snowflake id
	SnowFlake, err = goSnowFlake.NewIdWorker(1)
	if err != nil {
		panic(err)
	}

	//validate
	config := &validator.Config{TagName: "validate"}
	Validate = validator.New(config)

}

func LogInit() {
	Logs = logrus.New()
	Logs.Out = ioutil.Discard
	Logs.Level = logrus.DebugLevel
	Logs.Formatter = &logrus.TextFormatter{ForceColors: true, DisableTimestamp: true}

	hook, err := utils.CreateFileHook()
	if err != nil {
		panic(err)
	}
	Logs.Hooks.Add(hook)
}

func SqlInit() {
	var err error
	err = ConnectMysql()
	if err != nil {
		panic(err)
	}

	err = ConnectRedis()
	if err != nil {
		panic(err)
	}
}

func ConnectRedis() (err error) {
	RedisStore = utils.NewRedisStore(model.AppConfig.RedisHost, model.AppConfig.RedisAuth)
	return
}

func ConnectMysql() (err error) {
	MysqlEngine, err = xorm.NewEngine("mysql", model.AppConfig.MysqlDsn)
	if err != nil {
		return
	}
	MysqlEngine.Sync2(
		new(model.User),
		new(model.Journal),
		new(model.Comment),
	)
	MysqlEngine.ShowSQL(true)
	return
}

func GetSnowFlakeId() int64 {
	id, err := SnowFlake.NextId()
	if err != nil {
		panic(err)
	}
	return id
}
