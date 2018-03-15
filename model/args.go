package model

type SignUpArgs struct {
	Alias    string `json:"alias" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginArgs struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JournalArgs struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Public  string `json:"public" binding:"required"`
}
