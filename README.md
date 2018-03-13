# Journal


Journal

    -- config   //json配置表

    -- controller //纯业务逻辑

    -- doc      //文档

    -- model    //纯数据

    -- router  //路由, 中间件

    -- service //app功能支持

    -- utils  //工具




### gin log

```
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //log+console
```

### json
https://github.com/json-iterator/go

### validator
https://github.com/go-playground/validator

### redis
https://github.com/garyburd/redigo/