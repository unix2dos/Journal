package service

import (
	"Journal/model"
	"os"

	"github.com/gin-gonic/contrib/cache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"github.com/zheng-ji/goSnowFlake"
)

var (
	MysqlEngine *xorm.Engine
	RedisStore  *cache.RedisStore
	SnowFlake   *goSnowFlake.IdWorker
	Logs        *logrus.Entry
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
}

func LogInit() {

	file, err := os.OpenFile("logs/main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logs := logrus.New()
	logs.Out = file
	logs.Level = logrus.DebugLevel
	logs.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	Logs = logs.WithField("MODULE", "API")
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
	RedisStore = cache.NewRedisCache(model.AppConfig.RedisHost, model.AppConfig.RedisAuth, cache.DEFAULT)
	return
}

func ConnectMysql() (err error) {
	MysqlEngine, err = xorm.NewEngine("mysql", model.AppConfig.MysqlDsn)
	if err != nil {
		return
	}
	MysqlEngine.Sync2(new(model.User))
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
