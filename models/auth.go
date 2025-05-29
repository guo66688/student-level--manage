package models

type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"123456"`
}
