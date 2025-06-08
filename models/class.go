// models\class.go
package models

import "gorm.io/gorm"

// Class 班级信息
type Class struct {
	gorm.Model
	Name     string    `gorm:"unique" json:"name"`
	Students []Student `gorm:"foreignKey:ClassID" json:"students"`
}
