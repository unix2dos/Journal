package controller

import (
	"Journal/model"
	"Journal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequestLog(c *gin.Context) {

	requestDump, _ := httputil.DumpRequest(c.Request, true)
	start := time.Now()
	c.Next()

	latency := time.Since(start)
	v, _ := c.Get("data")
	bytes, _ := json.Marshal(v)
	service.Logs.Infof("request=| %v \nresponse=| status=%v ip=%v duration=%v data=%v",
		string(requestDump),
		c.Writer.Status(),
		c.ClientIP(),
		latency,
		string(bytes),
	)
}

func SessionFilter(c *gin.Context) {
	//siginup和login不经过
	if c.Request.RequestURI == "/signup" || c.Request.RequestURI == "/login" {
		return
	}
	data := GetData(c)

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
	data.Data["user"] = user //TODO: 如果写到这里, user信息变了, 可能还是老数据, 如果写到commonReturn里呢??
}

func CommonReturn(c *gin.Context) {
	c.Next()
	data := GetData(c)
	data.Msg = model.GetDataMsg(data.Ret)
	data.Data["ts"] = fmt.Sprintf("%d", time.Now().Unix())

	if data.Ret == model.Success {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusInternalServerError, data)
	}
}
