package models

import "gorm.io/gorm"

// Student 学生信息
type Student struct {
	gorm.Model
	StudentNo string `gorm:"type:varchar(255);uniqueIndex" json:"student_no" example:"20230001"`
	Name      string `gorm:"type:varchar(255)" json:"name" example:"张三"`
	Gender    string `gorm:"type:varchar(50)" json:"gender" example:"男"`
	BirthDate string `gorm:"type:varchar(255)" json:"birth_date" example:"2005-09-01"`

	ClassID uint  `json:"class_id"`                        // 外键字段
	Class   Class `gorm:"foreignKey:ClassID" json:"class"` // 自动加载关联班级
}
