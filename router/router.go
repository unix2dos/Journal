package router

import (
	"Journal/controller"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(controller.CommonReturn)

	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)
}
