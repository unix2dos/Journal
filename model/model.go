package model

//全局数据
var (
	AppConfig = &config{}
)

//redis key
var (
	RedisKeyPrefix  string = "journal"
	RedisKeyUser    string = RedisKeyPrefix + ":user:%d"
	RedisKeyJournal string = RedisKeyPrefix + ":journal:%d"
	RedisKeyComment string = RedisKeyPrefix + ":comment:%d"
)
