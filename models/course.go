// models\course.go
package models

import "gorm.io/gorm"

// Course 课程信息
type Course struct {
	gorm.Model
	CourseName string  `json:"course_name" example:"数学"`
	CourseCode string  `json:"course_code" gorm:"unique" example:"MATH101"`
	Credit     float64 `json:"credit" example:"3"`
	TeacherID  uint    `json:"teacher_id" example:"10"`
}
