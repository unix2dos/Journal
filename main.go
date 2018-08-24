package main

import (
	"github.com/gin-gonic/gin"

	"github.com/unix2dos/Journal/controller"
	"github.com/unix2dos/Journal/model"
	"github.com/unix2dos/Journal/service"
	"github.com/unix2dos/Journal/utils"
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
