tags:

#### 1. 项目可以用 jsoniter, 比原生的快
	go build -tags=jsoniter .
	
#### 2. router

	+ /user/name 
	
	```
		r.GET("/user/:name/*action", func(c *gin.Context) {

		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
	```
	+ /welcome?name=levon&sex=boy

	```
		r.GET("/welcome", func(c *gin.Context) {
		name := c.DefaultQuery("name", "name")
		sex := c.Query("sex")
		c.String(http.StatusOK, name+sex)
	})
	```
	
	+ post form data

	```
		r.POST("/formpost", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "liuwei")
		c.JSON(http.StatusOK, gin.H{"status": "posted", "message": message, "nick": nick})
	})
	```
	
	+ upload file

	```
		r.POST("upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		c.SaveUploadedFile(file, "a.png")
		c.String(http.StatusOK, file.Filename)
	})
	```
	
	+ group router

	```
	v1 := r.Group("/v1")
	{
		v1.POST("/get", func(c *gin.Context) {
			c.String(http.StatusOK, "v1")
		})
	}

	v2 := r.Group("/v2")
	{
		v2.POST("/get", func(c *gin.Context) {
			c.String(http.StatusOK, "v2")
		})
	}
	```

#### 3. middleware

	+ use

	```
	r := gin.New()
	r.Use(gin.Logger())
	```

	+ group

	```
		authorized := r.Group("/")
	authorized.Use(func(c *gin.Context) {
		log.Println("custom middle")//调用了这个
	})
	{
		authorized.GET("/get", func(c *gin.Context) {
			c.String(http.StatusOK, "1")
		})

		// nested group	然是nest,但是访问不是localhost:8080/testing/get
		testing := authorized.Group("testing")
		testing.GET("/get", func(c *gin.Context) {
			c.String(http.StatusOK, "2")
		})
	}
	```
	
#### 4. write log

+ 日志写的是请求信息, 不是log.Println
+ gin.DefaultWriter还要在gin.New()之前

```
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //log+console

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		log.Println("haha")
		c.String(200, "pong")
	})
```


#### 5. model bind

1. 绑定依赖于 Content-Type header
+ Bind, BindJSON, BindQuery 如果出错是400

```
	r.POST("/login_json", func(c *gin.Context) {
		var json LOGIN
		c.BindJSON(&json)
	
		//如果上面出错了, 现在401错误也是400
		if json.Name == "liuwei" && json.Pass == "1" {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauth"})
		}
	})
	
	//{"name":"liuwei","pass":"12"}
```

+ ShouldBind, ShouldBindJSON, ShouldBindQuery 出错,开发者处理
+ BindQuery 只能绑定get, 不能绑定post数据 
+ Bind()是可以bind Query String和 post data


#### 6. return data

	c.String
	c.JSON
	c.XML
	c.YAML

#### 7. custom middleware

```
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("name", "liuwei")
		log.Print(t)

		c.Next()//此处前是请求前, 后面是请求后, 功能是相当于提前执行了handlers

		t1 := time.Since(t)
		log.Print(t1)

		status := c.Writer.Status()
		log.Println(status)
	}
}
```

#### 8. BasicAuth授权

```
	group := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":  "bar",
		"name": "liuwei",
	}))
	group.GET("/secrets", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		log.Println(user)
	})
```

#### 9. gin.Context Groutines

如果用groutines需要copy, 不copy什么后果?

```
	r.GET("/long_async", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		go func() {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)

			// note that you are using the copied context "cCp", IMPORTANT
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})
```







