package models

import "gorm.io/gorm"

// Class 班级信息
type Class struct {
	gorm.Model
	Name     string    `gorm:"unique" json:"name" example:"高一1班"` // 班级名称
	Students []Student `gorm:"foreignKey:ClassName;references:Name" json:"students"`
}
