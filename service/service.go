package service

import (
	"Journal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/zheng-ji/goSnowFlake"
)

var (
	MysqlEngine *xorm.Engine
	SnowFlake   *goSnowFlake.IdWorker
)

func SInit() {
	LogInit()
	SqlInit()

	var err error
	SnowFlake, err = goSnowFlake.NewIdWorker(1)
	if err != nil {
		panic(err)
	}
}

func LogInit() {
	//TODO: log
}

func SqlInit() {
	var err error
	err = ConnectMysql(model.AppConfig.MysqlDsn)
	if err != nil {
		panic(err)
	}
}

func ConnectMysql(conn string) (err error) {
	MysqlEngine, err = xorm.NewEngine("mysql", conn)
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
