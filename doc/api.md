# 统一返回
### RetCode
错误码

```
	Success   = 0
	ErrorArgs = iota //1
	ErrorValidate

	ErrorUserNotExist
	ErrorRepeatSignUp
	ErrorUserPassWord
	ErrorJournalNotExist
	ErrorLikeAlready
	ErrorLikeNotExist

	ErrorServe = -1
```
错误码描述

```
	Success:       "success",
	ErrorServe:    "serve error",
	ErrorArgs:     "error args",
	ErrorValidate: "err validate",

	ErrorUserNotExist:    "user not exist",
	ErrorRepeatSignUp:    "email has sign up",
	ErrorUserPassWord:    "email or password error",
	ErrorJournalNotExist: "journal not exist",
	ErrorLikeAlready:     "already liked",
	ErrorLikeNotExist:    "have not liked",
```

### 返回统一结构
所有return 结构均为：

```
	{
		'ret': RetCode,
		'msg': "对应RetCode",
		'data': {} （不同接口不一样）
	 }
```
# API 接口信息


####  /getinfo GET 获取用户信息,有cookie直接返回

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "time": "1521635430",
        "user": {
            "id": "196378450074800128",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```



####  /signup POST 注册

字段|说明
---|---
alias| 昵称
email| 邮箱
password| 密码

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "time": "1521635485",
        "user": {
            "id": "196430148772302848",
            "alias": "13",
            "email": "13@qq.com"
        }
    }
}
```


####  /login POST 登录

字段|说明
---|---
email | 邮箱
password| 密码

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "time": "1521635574",
        "user": {
            "id": "196378450074800128",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```


#### /journal/list GET 获取日志信息


```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "journals": [
            {
                "id": "196381018574295040",
                "title": "555",
                "content": "555",
                "public": "1",
                "user_id": 196378450074800128,
                "create_time": "1521623771",
                "update_time": "1521627520",
                "like_count": "2",
                "like_by_me": "0"
            },
            {
                "id": "196381079895019520",
                "title": "e",
                "content": "e",
                "public": "1",
                "user_id": 196378450074800128,
                "create_time": "1521623786",
                "update_time": "1521623786",
                "like_count": "3",
                "like_by_me": "1"
            }
        ],
        "time": "1521635766",
        "user": {
            "id": "196378450074800128",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```

#### journal/add POST 增加日志
字段|说明
---|---
title | 标题
content| 内容
public|是否公开 "1" "0"

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "journal": {
            "id": "196482399075307520",
            "title": "e",
            "content": "e",
            "public": "1",
            "user_id": 196474816105025536,
            "create_time": "1521647942",
            "update_time": "1521647942",
            "like_count": "0",
            "like_by_me": ""
        },
        "time": "1521647942",
        "user": {
            "id": "196474816105025536",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```
#### journal/update POST 修改日志
字段|说明
---|---
id|日志id
title | 标题
content| 内容
public|是否公开 "1" "0"

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "journal": {
            "id": "196482399075307520",
            "title": "555",
            "content": "555",
            "public": "0",
            "user_id": 196474816105025536,
            "create_time": "1521647942",
            "update_time": "1521648020",
            "like_count": "0",
            "like_by_me": "0"
        },
        "time": "1521648020",
        "user": {
            "id": "196474816105025536",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```


#### journal/delete POST 删除日志
字段|说明
---|---
id|日志id

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "time": "1521648084",
        "user": {
            "id": "196474816105025536",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```

#### journal/recommend GET 获取推荐日志

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "journals": [
            {
                "id": "196473841453633536",
                "title": "1->3",
                "content": "3",
                "public": "1",
                "user_id": 196473692539064320,
                "create_time": "1521645902",
                "update_time": "1521646074",
                "like_count": "1",
                "like_by_me": "0"
            },
            {
                "id": "196473821438414848",
                "title": "1->2",
                "content": "2",
                "public": "1",
                "user_id": 196473692539064320,
                "create_time": "1521645897",
                "update_time": "1521646061",
                "like_count": "1",
                "like_by_me": "0"
            }
        ],
        "time": "1521648183",
        "user": {
            "id": "196483390575218688",
            "alias": "13",
            "email": "13@qq.com"
        }
    }
}
```

#### /like/add POST 喜欢日志或喜欢评论

字段|说明
---|---
like_type|喜欢类型 "1"->journal "2"->comment
like_id|喜欢id

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "time": "1521648319",
        "user": {
            "id": "196483390575218688",
            "alias": "13",
            "email": "13@qq.com"
        }
    }
}
```


#### /like/delete POST 删除日志或删除评论

字段|说明
---|---
like_type|喜欢类型 "1"->journal "2"->comment
like_id|喜欢id

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "time": "1521648319",
        "user": {
            "id": "196483390575218688",
            "alias": "13",
            "email": "13@qq.com"
        }
    }
}
```

### /comment/list GET 获取单个日志评论列表

字段|说明
---|---
JournalId|日志id  例如: /comment/list?JournalId=xxxxxxx

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "comments": [
            {
                "id": "197193699036237824",
                "content": "i am a comment",
                "journal_id": "197193667373436928",
                "create_time": "1521817529",
                "update_time": "1521817529",
                "user_alias": "2"
            }
        ],
        "time": "1521817573",
        "user": {
            "id": "197188279152414720",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```


### /comment/add POST 增加评论

字段|说明
---|---
journal_id|日志id
content|评论内容
`reply_comment_id`| 回复评论的id, 默认是0

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "comment": {
            "id": "197193699036237824",
            "content": "i am a comment",
            "journal_id": "197193667373436928",
            "create_time": "1521817529",
            "update_time": "1521817529",
            "user_alias": "2"
        },
        "time": "1521817529",
        "user": {
            "id": "197188279152414720",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```

### /comment/update
字段|说明
---|---
comment_id|评论id
content|评论内容

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "comment": {
            "id": "197193699036237824",
            "content": "modify content",
            "journal_id": "197193667373436928",
            "create_time": "1521817529",
            "update_time": "1521817637",
            "user_alias": "2"
        },
        "time": "1521817637",
        "user": {
            "id": "197188279152414720",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```


### /comment/delete
字段|说明
---|---
comment_id|评论id

```
{
    "ret": 0,
    "msg": "success",
    "data": {
        "time": "1521817677",
        "user": {
            "id": "197188279152414720",
            "alias": "2",
            "email": "2@qq.com"
        }
    }
}
```
