结合“好鱼丸.pdf” 中的接口, 调整夏雨2018-01-24时的backend/app.py. 尚未测试。
server 地址暂定journalex.us:5000

# Codes as outlined in the current frontend/backend communication document

CODE_SUCCESS = '0'
CODE_ERR_EMAIL = '1001'
CODE_ERR_EMAIL_PASS_NULL = '1002'
CODE_ERR_EMAIL_PASS_WRONG = '1003'
CODE_ERR_EMAIL_TAKEN = '1004'
CODE_ERR_EMAIL_NOT_FOUND = '1005'
CODE_ERR_VERIFICATION = '1006'
CODE_ERR_VALIDATION = '1007'
CODE_ERR_NOTLOGGEDIN = '1008'
CODE_ERR_USRNAMETAKEN = '1009'
CODE_ERR_BADINVITECODE = '1010'


CODE_ERR_DB = '2001'
CODE_ERR_INTERNAL = '2002'
CODE_ERR_ALREADYLOGIN = '2003'


# Only the msg associated with the error codes that are currently used
MSG_SUCCESS = 'Request Succeeded'
MSG_ERR_EMAIL_PASS_NULL = 'Username/Password cannot be empty'
MSG_ERR_EMAIL_PASS_WRONG = 'Username/Password is incorrect'
MSG_ERR_VALIDATION = 'Did not pass validator tests.'
MSG_ERR_NOTLOGGEDIN = 'Not logged in.'
MSG_ERR_ALREADYLOGIN = 'Internal Error: an already logged in session is directed to sign-up page'
MSG_ERR_USRNAMETAKEN = 'Username already taken'
MSG_ERR_BADINVITECODE = 'Bad Invitation Code'


所有return 结构均为：
{
            'ret': CODE_******,
            'msg': MSG_*******,
            'data': ******** （不同functions 不一样）
 }




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




三. /getinfo
Method: Get 
Success:
接口返回数据和login一样


四. /logout
Method: Get
Success:
Return to /login page



五. /signup
Method: Post (username, password,alias)
Success:
{
              'ret': CODE_SUCCESS,
              'msg': MSG_SUCCESS,
              'data': data_dict
            }
data_dict 同二。


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






