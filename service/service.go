package service

import (
	"Journal/model"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	MysqlEngine *xorm.Engine
	Log         *logrus.Entry
)

func SInit() {
	LogInit()
	SqlInit()
}

func LogInit() {
	//TODO: log
}

func SqlInit() {
	var err error
	err = ConnectDB(model.AppConfig.MysqlDsn)
	if err != nil {
		panic(err)
	}
}

func ConnectDB(conn string) (err error) {
	MysqlEngine, err := xorm.NewEngine("mysql", conn)
	if err != nil {
		return
	}
	MysqlEngine.Sync(new(model.User))
	return
}
