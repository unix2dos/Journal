package main

import (
	"os"

	"io"

	"Journal/router"

	"github.com/gin-gonic/gin"
)

func main() {

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //log+console

	r := gin.New()
	router.Route(r)

	r.Run()
}
