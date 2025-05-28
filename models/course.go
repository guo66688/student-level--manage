package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	CourseName string
	CourseCode string `gorm:"unique"` // 课程代码唯一
	Credit     float64
	TeacherID  uint // 外键，不设置约束
}
