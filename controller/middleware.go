package controller

import (
	"Journal/model"
	"Journal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"Journal/utils"

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

	if utils.StringContains(c.Request.RequestURI, NotSessionFilter) {
		return
	}

	data := GetData(c)
	userId, ok := SessionGet(c)
	//没有cookie
	if !ok {
		data.Ret = model.ErrorNotLogin
		c.Abort()
		return
	}

	//有cookie
	user, exist, err := userService.GetUserById(userId)
	if err != nil {
		c.Abort()
		data.Ret = model.ErrorServe
	}

	//有cookie没有用户
	if !exist {
		c.Abort()
		data.Ret = model.ErrorNotLogin
		return
	}

	//有cookie, 有用户
	c.Set("uid", user.Id)

	//data.Data["user"] = user //TODO: 如果写到这里, user信息变了, 可能还是老数据, 如果写到commonReturn里呢??
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
