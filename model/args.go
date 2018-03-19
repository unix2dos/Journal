package model

type SignUpArgs struct {
	Alias    string `json:"alias" binding:"required" validate:"gt=0,lte=50"`
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"gte=6,lte=50"`
}

type LoginArgs struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"gte=6,lte=50"`
}

type JournalAddArgs struct {
	Title   string `json:"title" binding:"required" validate:"gt=0"`
	Content string `json:"content" binding:"required" validate:"gt=0"`
	Public  string `json:"public" binding:"required" validate:"eq=1"`
}

type JournalUpdateArgs struct {
	Id      string `json:"id" binding:"required" validate:"gt=0"`
	Title   string `json:"title" binding:"required" validate:"gt=0"`
	Content string `json:"content" binding:"required" validate:"gt=0"`
	Public  string `json:"public" binding:"required" validate:"eq=1"`
}

type JournalDeleteArgs struct {
	Id string `json:"id" binding:"required" validate:"gt=0"`
}
