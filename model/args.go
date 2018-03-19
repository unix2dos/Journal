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
