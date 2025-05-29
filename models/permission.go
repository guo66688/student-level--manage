// models/permission.go
package models

import "gorm.io/gorm"

// 权限定义
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"` // 权限标识，如 "view_chart"
	Description string         `json:"description"`                 // 权限说明
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// 角色定义
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"` // 角色标识，如 "admin"
	Description string         `json:"description"`                 // 角色说明
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// 角色权限关联表
type RolePermission struct {
	ID           uint `gorm:"primaryKey" json:"id"`
	RoleID       uint `gorm:"index;not null" json:"role_id"`
	PermissionID uint `gorm:"index;not null" json:"permission_id"`
}

// 用户角色关联表
type UserRole struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index;not null" json:"user_id"`
	RoleID uint `gorm:"index;not null" json:"role_id"`
}
