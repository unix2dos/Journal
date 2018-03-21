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

#TODO: 留言系统

<!--留言系统

十二.comment/add
Method: POST (journal_id [int], content)
journal_id: 从journal/list 或 journal/recommend 里获取
content {
	“text”: 字符串
}
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': comment id
 }

十三.comment/addreply （回应别人的留言）
Method: POST (journal_id [int], content, target_user_id[int])
journal_id: 从comment/get 里获取 
target_user_id: 从comment/get 里获取
content: 同十二
 
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': comment id
 }

十四. comment/update
Method: Post (comment_id [int], content)
comment_id: 每个comment(留言）都有一个unique id, 从comment/get 里获取 
content: 同十二

Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': {}
 }


十五. comment/get/jid=<jid> 列出对应journal的所有相关留言。如果是journal 作者，列出所有留言。如果是旁人（我），列出journal 作者公开留言(journal作者使用comment/add), 和journal作者回复我的留言（journal 作者使用comment/addreply, 其中target_user_id 是我的user_id),以及我对该journal的所有留言
Method: Get
jid: journal id, 从journal/list 或 journal/recommend 里获取

Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': list of comment
 }
comment {
'content': 同十二    	'id': 留言id    	'journal_id': 留言所属journal 的id，    	'misc': 暂时没用，    	'target_user_id': {None (用comment/add 提交的留言) 或者是user id (用comment/addreply 提交的留言)},    	'user_id': 留言作者,
'timestamp_create',
'timestamp_update'}
}


十六. comment/delete
Method: Post (comment_id)
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': {}
 }



点赞系统

十七. like/get/jid=<jid>&cid=<cid>
Method: Get
jid: journal_id 从journal list 获取
cid: comment_id 从 comment/get 获取
jid 和 cid 其中一个是-1. 比如我想点赞journal 3, jid=3, cid = -1。 如果我想点赞journal 3 下面 comment 2, jid=-1, cid = 2。
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': {‘like_count': 这的journal/comment 的总赞数，
’i_like': 我有没有点赞， 0 是没有， 1是有,
'likeList‘: list of like 
 }
like {comment_id, journal_id, id, user_id}

十八. like/add
method: Post(journal_id, comment_id) 原理同上
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': {}
 }

十九.like/delete
method: Post(like_id) 每个like entry 都有自己的id, 可以从like/get 获取
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': {}
 }

二十.archive/get (list of 我点赞了的journals）
method: Get
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': list of journals
 }
journal 结构同七

注，点赞系统的一个重要环节在于给别人journal点赞了才能进行留言。我设想这个功能由前端实现。 前端可以通过like/get 判断journals 是不是被我点赞了，从而决定显示 留言功能。-->






