package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"` // 登录名唯一
	Password string // 明文测试环境，建议上线前加密
	Role     string // teacher / admin
}
