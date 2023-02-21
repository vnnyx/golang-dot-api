package web

type UserCreateRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	Handphone            string `json:"handphone"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type UserResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Handphone string `json:"handphone"`
}

type UserResponseWithLastTransaction struct {
	UserID      string              `json:"user_id"`
	Username    string              `json:"username"`
	Email       string              `json:"email"`
	Handphone   string              `json:"handphone"`
	Transaction TransactionResponse `json:"last_transaction"`
}

type UserUpdateProfileRequest struct {
	UserID    string
	Username  string `json:"username"`
	Email     string `json:"email"`
	Handphone string `json:"handphone"`
}

type UserUpdatePasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
