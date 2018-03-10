### TODO
1. redis type
2. read all

### RetCode

```
	Success    = 0
	ErrorServe = iota
	ErrorArgs
	ErrorNotLogin
	ErrorRepeatSignUp
	ErrorUserPassWord
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




## API

###  /getinfo 获取用户信息 GET



###  /signup 注册 POST

字段|说明|
---|---|---
alias| 昵称
email| 邮箱
password| 密码

###  /login 登录 POST

字段|说明|
---|---|---
email | 邮箱
password| 密码



二. /login
Method: POST (username, password)
Success:
{
            'ret': CODE_SUCCESS,
            'msg': MSG_SUCCESS,
            'data': data_dict
          }

Data_dict: 
{
    'basic': {},
    'userinfo': {
      'id': int,
      'username': email,
      'alias': string,
      'likes': int,
      ‘journal_count’: int,
      ‘avatar’: string indicating url to avatar
    },
    'journals': list of journals (见下）, 
  }

journal:
{
“title”:string,
“tags”: string （’#tag1#tag2#tag3 ...’,  // 分隔符为#）,
“content”:{content},   // content也是一个json对象
“published”:0,  // 0 或 1，1为published
"timestamp_create": 字符串
"timestamp_update": 字符串

}
content为一个json对象，目前的结构为
{
‘text’:’today is a good day’, 
}
只用它来存放纯文本，所以只用’text’一个field。











七. /journal/list/public={-1,0,1} 
-1: 所有；0: 私有； 1：公开
Method: Get
Success:
{
      'ret': CODE_SUCCESS,
      'msg': MSG_SUCCESS,
      'data': journallist
    }
Journallist: list of journal
Journal:
{
‘id’:int,
‘title’:’aaa’,
‘tags’:’tag1 tag2 ...’,  // 分隔符为空格
‘content’:{content},   // content也是一个json对象
‘time’:Datetime obj,
‘published’:0,  // 0 或 1，1为published
}

content为一个json对象，目前的结构为
{
‘text’:’today is a good day’, 
}
只用它来存放纯文本，所以只用’text’一个field。


八. journal/check/jid={从 journal/list 里来的journal id}, 如果用别的id 会有internal server error
Method: get
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': journal (同 二里面的定义）
      }

九. journal/add
Method: Post (tags, title, content, published) [content 为文本，published 为{0，1}； 0是私有，1是公开]
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': journal.id
      }


十.journal/update
Method: Post (id, {tags, title, content, published}), id 必要
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': journal.id
      }

十一.journal/recommend
Method: Get
Success:
{
        'ret': CODE_SUCCESS,
        'msg': MSG_SUCCESS,
        'data': list of journals
 }

留言系统

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

注，点赞系统的一个重要环节在于给别人journal点赞了才能进行留言。我设想这个功能由前端实现。 前端可以通过like/get 判断journals 是不是被我点赞了，从而决定显示 留言功能。






