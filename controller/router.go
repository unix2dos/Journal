package controller

import (
	"Journal/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) {
	store, err := sessions.NewRedisStore(100, "tcp", model.AppConfig.RedisHost, model.AppConfig.RedisAuth, []byte("secret"))
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions("journal", store))

	//中间件顺序不要变
	r.Use(gin.Logger()) //默认请求日志
	r.Use(gin.Recovery())
	r.Use(RequestLog)
	r.Use(CommonReturn)
	r.Use(SessionFilter)

	r.GET("/getinfo", GetInfo)
	r.POST("/signup", Signup)
	r.POST("/login", Login)

	journal := r.Group("/journal")
	{
		journal.GET("/list", JournalList)
		journal.POST("/add", JournalAdd)
		journal.POST("/update", JournalUpdate)
		journal.POST("/delete", JournalDel)

		journal.GET("/recommend", JournalRecommend)
	}

	comment := r.Group("/comment")
	{
		comment.GET("/list", CommentList)
		comment.POST("/add", CommentAdd)
		comment.POST("/update", CommentUpdate)
		comment.POST("/delete", CommentDelete)
	}

	like := r.Group("/like")
	{
		like.POST("/add", LikeAdd)
		like.POST("/delete", LikeDelete)
	}

	archive := r.Group("/archive")
	{
		archive.GET("/get", ArchiveGet)
	}

}
