package model

type SignUpArgs struct {
	Alias    string `json:"alias" binding:"required" validate:"lte=50"`
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"gte=6,lte=50"`
}

type LoginArgs struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"gte=6,lte=50"`
}

type JournalAddArgs struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Public  string `json:"public" binding:"required" validate:"len=1"`
}

type JournalUpdateArgs struct {
	Id      string `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Public  string `json:"public" binding:"required" validate:"len=1"`
}

type JournalDeleteArgs struct {
	Id string `json:"id" binding:"required"`
}

type LikeAddArgs struct {
	//"1"	journal
	//"2"	comment
	LikeType string `json:"like_type" binding:"required" validate:"eq=1|eq=2"`
	LikeId   int64  `json:"like_id,string" binding:"required"`
}

type LikeDelArgs struct {
	LikeType string `json:"like_type" binding:"required" validate:"eq=1|eq=2"`
	LikeId   int64  `json:"like_id,string" binding:"required"`
}

type CommentListArgs struct {
	JournalId int64 `json:"journal_id,string" binding:"required"`
}

type CommentAddArgs struct {
	JournalId      int64  `json:"journal_id,string" binding:"required"`
	Comment        string `json:"comment" binding:"required"`
	ReplyCommentId int64  `json:"reply_comment_id,string"`
}

type CommentUpdateArgs struct {
	CommentId int64 `json:"comment_id,string"`
}

type CommentDeleteArgs struct {
	CommentId int64 `json:"comment_id,string"`
}
