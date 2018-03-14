package main

import (
	"Journal/controller"
	"Journal/model"
	"Journal/service"
	"Journal/utils"

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
	r := gin.New()
	controller.Route(r)
	r.Run(model.AppConfig.RunPort)
}
