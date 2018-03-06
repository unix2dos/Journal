

### gin log

```
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //log+console
```

### json
https://github.com/json-iterator/go

### validator 

https://github.com/go-playground/validator