package models

import "gorm.io/gorm"

type Score struct {
	gorm.Model
	StudentID uint
	CourseID  uint
	Score     float64
	ExamDate  string
}
