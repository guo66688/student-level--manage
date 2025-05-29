package models

import "gorm.io/gorm"

// User 用户信息
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);uniqueIndex;not null" json:"username"` // ✅ 添加 type:varchar(255)
	Password string `gorm:"type:varchar(255);not null" json:"password"`             // ✅ 添加 type:varchar(255)
	Nickname string `json:"nickname"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}
