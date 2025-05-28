// models/score.go
package models

import "gorm.io/gorm"

type Score struct {
	gorm.Model
	StudentID uint    `json:"student_id"`
	CourseID  uint    `json:"course_id"`
	Score     float64 `json:"score"`
	ExamDate  string  `json:"exam_date"`
}
