package models

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Name     string    `gorm:"unique"` // 班级名称唯一
	Students []Student `gorm:"foreignKey:ClassName;references:Name"`
}
