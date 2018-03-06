package main

import (
	"Journal/router"

	"Journal/model"
	"Journal/utils"

	"github.com/gin-gonic/gin"
)

func init() {

	//configs
	configs := make(map[string]interface{})
	configs["./conf/config_test.json"] = model.AppConfig
	utils.ParseConfigs(configs)
}

func main() {

	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //log+console

	r := gin.New()
	router.Route(r)
	r.Run(model.AppConfig.RunPort)
}
