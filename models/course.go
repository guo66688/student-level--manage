package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	CourseName string  `json:"course_name"`
	CourseCode string  `json:"course_code" gorm:"unique"`
	Credit     float64 `json:"credit"`
	TeacherID  uint    `json:"teacher_id"`
}
