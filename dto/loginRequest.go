package dto

type LoginRequest struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}
