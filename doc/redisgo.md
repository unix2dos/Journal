### Executing Commands

支持redis所有command

```
Do(commandName string, args ...interface{}) (reply interface{}, err error)
```



### Args

Args is a helper for constructing command arguments from structured values.


```
var p1, p2 struct {
    Title  string `redis:"title"`
    Author string `redis:"author"`
    Body   string `redis:"body"`
}

p1.Title = "Example"
p1.Author = "Gary"
p1.Body = "Hello"

if _, err := c.Do("HMSET", redis.Args{}.Add("id1").AddFlat(&p1)...); err != nil {
    fmt.Println(err)
    return
}

m := map[string]string{
    "title":  "Example2",
    "author": "Steve",
    "body":   "Map",
}

if _, err := c.Do("HMSET", redis.Args{}.Add("id2").AddFlat(m)...); err != nil {
    fmt.Println(err)
    return
}

for _, id := range []string{"id1", "id2"} {

    v, err := redis.Values(c.Do("HGETALL", id))
    if err != nil {
        fmt.Println(err)
        return
    }

    if err := redis.ScanStruct(v, &p2); err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("%+v\n", p2)
}
```