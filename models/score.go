// models/score.go
package models

import "gorm.io/gorm"

// Score 学生成绩
type Score struct {
	gorm.Model
	StudentID uint    `json:"student_id" gorm:"index:idx_score,unique" example:"1"`
	CourseID  uint    `json:"course_id"  gorm:"index:idx_score,unique" example:"101"`
	ExamDate  string  `json:"exam_date"  gorm:"index:idx_score,unique" example:"2025-05-01"`
	Score     float64 `json:"score"      example:"89.5"`
}
