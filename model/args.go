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
	LikeType string `json:"like_type" binding:"required"` //"1"->journal  "2"->comment
	LikeId   int64  `json:"like_id" binding:"required"`
}

type LikeDelArgs struct {
	LikeType string `json:"like_type" binding:"required"`
	LikeId   int64  `json:"like_id" binding:"required"`
}
