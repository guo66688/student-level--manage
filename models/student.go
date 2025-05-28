package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	StudentNo string `gorm:"type:varchar(255);uniqueIndex"`
	Name      string `gorm:"type:varchar(255)"`
	Gender    string `gorm:"type:varchar(50)"`
	BirthDate string `gorm:"type:varchar(255)"`
	ClassName string `gorm:"type:varchar(255)"`
}
