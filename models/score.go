// models/score.go
package models

import "gorm.io/gorm"

// Score 学生成绩
type Score struct {
	gorm.Model
	StudentID uint    `json:"student_id" example:"1"`
	CourseID  uint    `json:"course_id" example:"101"`
	Score     float64 `json:"score" example:"89.5"`
	ExamDate  string  `json:"exam_date" example:"2025-05-01"`
}
