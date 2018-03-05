package main

import (
	"os"

	"io"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //log+console

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		log.Println("haha")
		c.String(200, "pong")
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
