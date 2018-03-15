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
		service.Logs.Errorf("ErrorNotLogin1 userId=%d", userId)
		data.Ret = model.ErrorNotLogin
		c.Abort()
		return
	}

	//有cookie
	user, exist, err := userService.GetUserById(userId)
	if err != nil {
		c.Abort()
		service.Logs.Errorf("ErrorServe userId=%d err=%v", userId, err)
		data.Ret = model.ErrorServe
	}

	//有cookie没有用户
	if !exist {
		c.Abort()
		service.Logs.Errorf("ErrorNotLogin2 userId=%d", userId)
		data.Ret = model.ErrorNotLogin
		return
	}

	//有cookie, 有用户
	c.Set("uid", user.Id)
}

func CommonReturn(c *gin.Context) {
	c.Next()

	if c.Writer.Status() == 404 || c.Writer.Status() == 405 {
		c.Abort()
		return
	}

	data := GetData(c)
	data.Msg = model.GetDataMsg(data.Ret)
	data.Data["time"] = fmt.Sprintf("%d", time.Now().Unix())

	if data.Ret == model.Success {

		uid, _ := c.Get("uid")
		user, exist, err := userService.GetUserById(uid.(int64))
		if err == nil && exist {
			data.Data["user"] = user
		} //TODO: 应该公共出去?

		c.JSON(http.StatusOK, data)

	} else if data.Ret == model.ErrorServe {

		c.JSON(http.StatusInternalServerError, data)

	} else {

		c.JSON(http.StatusOK, data)
	}

}
