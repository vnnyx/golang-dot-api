package web

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
}
