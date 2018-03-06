package service

import (
	"Journal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	MysqlEngine *xorm.Engine
)

func SInit() {
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
