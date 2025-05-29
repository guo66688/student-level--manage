package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"` // 登录名
	Password string `gorm:"not null"`             // 明文或加密密码
	Nickname string // 显示昵称
	// Role     string // teacher / admin，可弃用，用 RBAC 替代
	IsActive bool `gorm:"default:true"` // 是否启用
}
