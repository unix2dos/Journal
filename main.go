package main

import (
	"os"

	"io"

	"net/http"

	"github.com/gin-gonic/gin"
)

type LOGIN struct {
	Name string `form:"name" json:"name" binding:"required"`
	Pass string `form:"pass" json:"pass" binding:"required"`
}

func main() {

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //log+console

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/login_json", func(c *gin.Context) {
		var json LOGIN
		c.BindJSON(&json)

		if json.Name == "liuwei" && json.Pass == "1" {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauth"})
		}
	})

	r.POST("/login_form", func(c *gin.Context) {

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
