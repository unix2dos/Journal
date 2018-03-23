

## 创建xorm

```
 	var err error
   engine, err = xorm.NewEngine("mysql", "root:123@/test?charset=utf8")
```

+ engine是goroutine安全的
+ 创建完成engine之后，并没有立即连接数据库，此时可以通过engine.Ping()来进行数据库的连接测试是否可以连接到数据库。

## 日志

+ engine.ShowSQL(true)，则会在控制台打印出生成的SQL语句；
+ engine.Logger().SetLevel(core.LOG_DEBUG)，则会在控制台打印调试及以上的信息；

存到本地文件

```
f, err := os.Create("sql.log")
if err != nil {
    println(err.Error())
    return
}
engine.SetLogger(xorm.NewSimpleLogger(f))
```


## 名称映射规则

+ 默认SnakeMapper 支持struct为驼峰式命名，表结构为下划线命名之间的转换
+ SameMapper 支持结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名


gotype|sql
---|---|
int, int8, int16, int32, uint, uint8, uint16, uint32|	Int
int64, uint64|	BigInt
float32| Float
float64| Double
[]uint8| Blob
array, slice, map except []uint8| Text
bool| Bool
string|Varchar(255)
time.Time| DateTime
struct| Text

+ 使用Table和Tag改变名称映射
	* 如果结构体拥有TableName() string的成员方法，那么此方法的返回值即是该结构体对应的数据库表名。
	* 通过engine.Table()方法可以改变struct对应的数据库表的名称
	* 通过sturct中field对应的Tag中使用xorm:"'column_name'"可以使该field对应的Column名称为指定名称。

+ 表名的优先级顺序如下：
	* engine.Table() 指定的临时表名优先级最高
	* TableName() string 其次
	* Mapper 自动映射的表名优先级最后
+ 字段名的优先级顺序如下：
	* 结构体tag指定的字段名优先级较高
	* Mapper 自动映射的表名优先级较低


## Column tag 属性

列属性|解释|
---|---
name|对应的字段的名称
pk|Primary Key(int32,int,int64,uint32,uint,uint64,string这7种Go的数据类型)
autoincr|
null 或 notnull|
unique或unique(uniquename)|是否是唯一，如不加括号则该字段不允许重复；如加上括号，则括号中为联合唯一索引的名字
index或index(indexname)|是否是索引，如不加括号则该字段自身为索引，如加上括号，则括号中为联合索引的名字
-|	这个Field将不进行字段映射
->|	这个Field将只写入到数据库而不从数据库读取
<-|	这个Field将只从数据库读取，而不写入到数据库
default 0或default(0)|设置默认值，紧跟的内容如果是Varchar等需要加上单引号
json|表示内容将先转成Json格式，然后存储到数据库中，数据库中的字段类型可以为Text或者二进制
created|这个Field将在Insert时自动赋值为当前时间
updated|这个Field将在Insert或Update时自动赋值为当前时间
deleted|这个Field将在Delete时设置为当前时间，并且当前记录不删除
version|这个Field将会在insert时默认为1，每次更新自动加1


##### 潜规则
+ 如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。
+ string类型默认映射为varchar(255)，如果需要不同的定义，可以在tag中自定义，如：varchar(1024)
+ 支持type MyString string等自定义的field，支持Slice, Map等field成员，这些成员默认存储为Text类型，并且默认将使用Json格式来序列化和反序列化。


## 操作

##### 同步表结构
+ Sync	

	```
	err := engine.Sync(new(User), new(Group))
	```
+ Sync2

Sync2对Sync进行了改进，目前推荐使用Sync2。

##### 插入数据

+ 插入表

```
user := new(User)
user.Name = "myname"
affected, err := engine.Insert(user)
// INSERT INTO user (name) values (?)
```

+ 插入不同表

```
user := new(User)
user.Name = "myname"
question := new(Question)
question.Content = "whywhywhwy?"
affected, err := engine.Insert(user, question)
```

+ 时间

Created可以让您在数据插入到数据库时自动将对应的字段设置为当前时间，需要在xorm标记中使用created标记，如下所示进行标记，对应的字段可以为time.Time或者自定义的time.Time或者int,int64等int类型。

```
type User struct {
    Id int64
    Name string
    CreatedAt time.Time `xorm:"created"`
}
```

engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")


##### 查询数据

+ Alias(string) 给Table设定一个别名

	```
	engine.Alias("o").Where("o.name = ?", name).Get(&order)
	```

+ And(string, …interface{})
+ Asc(…string)
+ Desc(…string)
+ ID(interface{}) 传入一个主键字段的值，作为查询条件

	```
	var user User
	engine.ID(1).Get(&user)
	// SELECT * FROM user Where id = 1
	```
+ Or(interface{}, …interface{})
+ OrderBy(string)
+ Select(string)
+ SQL(string, …interface{}) 执行指定的Sql语句，并把结果映射到结构体。有时，当选择内容或者条件比较复杂时，可以直接使用Sql

	```
	engine.SQL("select * from table").Find(&beans)
	```
+ Where(string, …interface{})
+ In(string, …interface{})


#### Get

1. 查询单条数据使用Get方法，在调用Get方法时需要传入一个对应结构体的指针，同时结构体中的非空field自动成为查询的条件和前面的方法条件组合在一起查询。

2. 返回的结果为两个参数，一个has为该条记录是否存在，第二个参数err为是否有错误。不管err是否为nil，has都有可能为true或者false。

+ 根据Where来获得单条数据：

	```
	user := new(User)
	has, err := engine.Where("name=?", "xlw").Get(user)
	```

+ 根据user结构体中已有的非空数据来获得单条数据：

	```
	user := &User{Id:1}
	has, err := engine.Get(user)
	```
	
#### find

查询多条数据使用Find方法，Find方法的第一个参数为slice的指针或Map指针，即为查询后返回的结果，第二个参数可选，为查询的条件struct的指针。

+ 传入Slice用于返回数据

	```
	everyone := make([]Userinfo, 0)
	err := engine.Find(&everyone)
	
	pEveryOne := make([]*Userinfo, 0)
	err := engine.Find(&pEveryOne)
	```
+ 也可以加入各种条件

	```
	users := make([]Userinfo, 0)
	err := engine.Where("age > ? or name = ?", 30, "xlw").Limit(20, 10).Find(&users)
	```
	
#### update

更新数据使用Update方法，Update方法的第一个参数为需要更新的内容，可以为一个结构体指针或者一个Map[string]interface{}类型。当传入的为结构体指针时，只有非空和0的field才会被作为更新的字段。当传入的为Map类型时，key为数据库Column的名字，value为要更新的内容。

+ Update会自动从user结构体中提取非0和非nil得值作为需要更新的内容，因此，如果需要更新一个值为0，则此种方法将无法实现
 
	```
	user := new(User)
	user.Name = "myname"
	affected, err := engine.Id(id).Update(user)
	```


+ 指定列更新

	```
	affected, err := engine.Id(id).Cols("age").Update(&user)
	```

+ 通过传入map[string]interface{}来进行更新，但这时需要额外指定更新到哪个表
	
	```
	affected, err := engine.Table(new(User)).Id(id).Update(map[string]interface{}{"age":0})
	```
	
#### Delete

删除数据Delete方法，参数为struct的指针并且成为查询条件。

```
user := new(User)
affected, err := engine.Id(id).Delete(user)
```


#### Query
也可以直接执行一个SQL查询，即Select命令。

```
sql := "select * from userinfo"
results, err := engine.Query(sql)
```

#### Exec
也可以直接执行一个SQL命令，即执行Insert， Update， Delete 等操作。此时不管数据库是何种类型，都可以使用 ` 和 ? 符号。

```
sql = "update `userinfo` set username=? where id=?"
res, err := engine.Exec(sql, "xiaolun", 1) 
```











