package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Journal/controller"
	"Journal/model"
	"Journal/service"
	"Journal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ClientRequestLog(c *gin.Context) {
	var header string
	for k, v := range c.Request.Header {
		header += k
		for _, vv := range v {
			header += vv + " "
		}
	}

	////bytes, _ := ioutil.ReadAll(c.Request.Body)
	////log.Println(string(bytes))//TODO: fuck
	service.Logs.Debug(c.Request.Method, " ", c.Request.RequestURI, " ", header)

}

func ClientResponseLog(c *gin.Context) {
	start := time.Now()
	c.Next()

	latency := time.Since(start)
	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	statusColor := utils.ColorForStatus(statusCode)
	methodColor := utils.ColorForMethod(method)

	v, _ := c.Get("data")
	bytes, _ := json.Marshal(v)

	service.Logs.Debug(statusColor, statusCode, " ",
		utils.Reset, latency, " ", clientIP, " ",
		methodColor, method, " ",
		utils.Reset, c.Request.RequestURI, " ", string(bytes),
	)
}

func SessionFilter(c *gin.Context) {
	//siginup和login不经过
	if c.Request.RequestURI == "/signup" || c.Request.RequestURI == "/login" {
		return
	}
	data := controller.GetData(c)

	session := sessions.Default(c)
	userId, ok := session.Get("uid").(int64)
	if !ok {
		data.Ret = model.ErrorNotLogin
		c.Abort()
		return
	}

	user := new(model.User)
	exist, err := service.MysqlEngine.Id(userId).Get(user)
	if err != nil {
		data.Ret = model.ErrorServe
		c.Abort()
		return
	}

	if !exist { //cookie无效
		data.Ret = model.ErrorNotLogin
		c.Abort()
		return
	}
	c.Set("uid", user.Id)
	data.Data["user"] = user //TODO: 如果写到这里, user信息变了, 可能还是老数据
}

func CommonReturn(c *gin.Context) {
	c.Next()
	data := controller.GetData(c)
	data.Msg = model.GetDataMsg(data.Ret)
	data.Data["ts"] = fmt.Sprintf("%d", time.Now().Unix())

	if data.Ret == model.Success {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusInternalServerError, data)
	}
}
