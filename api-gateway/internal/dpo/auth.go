package dpo

type Regisration struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
