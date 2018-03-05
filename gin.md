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
