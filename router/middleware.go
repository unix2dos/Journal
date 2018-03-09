package router

import (
	"net/http"

	"Journal/controller"
	"Journal/model"
	"Journal/service"

	"fmt"
	"time"

	"log"

	"io/ioutil"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ClientRequestLog(c *gin.Context) {

	log.Println("-----------------")
	log.Println(c.Request.Method + " " + c.Request.RequestURI)
	log.Println("\n")
	for k, v := range c.Request.Header {
		log.Println(k, v)
	}
	log.Println("\n")
	bytes, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(bytes))
	log.Println("-----------------")
}

func ClientResponseLog(c *gin.Context) {
	c.Next()

	log.Println("+++++++++++++++++")
	log.Println(c.Request.Method + " " + c.Request.RequestURI)

	for k, v := range c.Writer.Header() {
		log.Println(k, v)
	}

	////c.Request.Body
	//bytes, _ := ioutil.ReadAll(c.Writer)
	//log.Println(string(bytes))
	log.Println("+++++++++++++++++")
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

	if !exist {
		//TODO: 无效cookie
	}

	data.Data["user"] = user
	data.Data["ts"] = fmt.Sprintf("%d", time.Now().Unix())
	c.Set("uid", user.Id)

}

func CommonReturn(c *gin.Context) {
	c.Next()
	log.Println("CommonReturn")
	data := controller.GetData(c)
	data.Msg = model.GetDataMsg(data.Ret)

	if data.Ret == model.Success {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusInternalServerError, data)
	}
}
