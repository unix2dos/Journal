package main

import (
	"Journal/router"

	"Journal/model"
	"Journal/utils"

	"Journal/service"

	"github.com/gin-gonic/gin"
)

func init() {

	//configs
	configs := make(map[string]interface{})
	configs["./conf/config_test.json"] = model.AppConfig
	utils.ParseConfigs(configs)

	//init
	service.SInit()
}

func main() {

	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //log+console

	r := gin.New()
	router.Route(r)
	r.Run(model.AppConfig.RunPort)
}
