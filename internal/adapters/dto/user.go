package dto

type SignupPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	ID       uint   `json:"id"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
