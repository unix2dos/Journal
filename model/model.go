package model

//全局数据
var (
	AppConfig = &config{}
)

//redis key
var (
	RedisKeyPrefix string = "journal"
	RedisKeyUser   string = RedisKeyPrefix + ":user:%d"
)
