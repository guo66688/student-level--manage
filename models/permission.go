// models/permission.go
package models

import "gorm.io/gorm"

// Permission 权限
// @Description 权限定义，例如 "student:view"
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id" example:"1"`
	Name        string         `gorm:"unique;not null" json:"name" example:"student:view"` // 权限标识
	Description string         `json:"description" example:"查看学生列表"`                       // 权限说明
	CreatedAt   int64          `json:"created_at" example:"1716940800"`
	UpdatedAt   int64          `json:"updated_at" example:"1716940800"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Role 角色
// @Description 角色定义，例如 "admin"
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id" example:"1"`
	Name        string         `gorm:"unique;not null" json:"name" example:"admin"` // 角色标识
	Description string         `json:"description" example:"系统管理员"`                 // 角色说明
	CreatedAt   int64          `json:"created_at" example:"1716940800"`
	UpdatedAt   int64          `json:"updated_at" example:"1716940800"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RolePermission 角色权限绑定
// @Description 角色与权限的绑定关系
type RolePermission struct {
	ID           uint `gorm:"primaryKey" json:"id" example:"1"`
	RoleID       uint `gorm:"index;not null" json:"role_id" example:"1"`
	PermissionID uint `gorm:"index;not null" json:"permission_id" example:"2"`
}

// UserRole 用户角色绑定
// @Description 用户与角色的绑定关系
type UserRole struct {
	ID     uint `gorm:"primaryKey" json:"id" example:"1"`
	UserID uint `gorm:"index;not null" json:"user_id" example:"3"`
	RoleID uint `gorm:"index;not null" json:"role_id" example:"1"`
}
