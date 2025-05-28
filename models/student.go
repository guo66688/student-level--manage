package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	StudentNo string `gorm:"uniqueIndex"` // 学号唯一
	Name      string
	Gender    string
	BirthDate string
	ClassName string
}
