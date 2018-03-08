package router

import (
	"Journal/controller"

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
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CommonReturn)
	r.Use(SessionFilter)

	r.GET("/getinfo", controller.GetInfo)
	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)
}
