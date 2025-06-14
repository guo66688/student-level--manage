// models\auth.go
package models

// LoginRequest 登录请求体
type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"123456"`
}
