// models\doc_analytics.go
package models

// CourseStats 用于课程统计响应
type CourseStats struct {
	CourseName string  `json:"course_name"`
	AvgScore   float64 `json:"avg_score"`
	MaxScore   float64 `json:"max_score"`
	MinScore   float64 `json:"min_score"`
}

// Result 用于月度平均成绩
type Result struct {
	Month string  `json:"month" example:"2024-06"`
	Avg   float64 `json:"avg" example:"82.7"`
}

// Pass 用于课程通过率
type Pass struct {
	CourseID   uint    `json:"course_id"` // ✅ 添加这一行
	CourseName string  `json:"course_name"`
	Passed     int64   `json:"passed"`
	Total      int64   `json:"total"`
	Rate       float64 `json:"rate"`
}

// ScoreRankRow 用于成绩排行榜响应
type ScoreRankRow struct {
	StudentName string  `json:"student_name" example:"张三"`
	AvgScore    float64 `json:"avg_score" example:"87.5"`
}

// ClassAvgDoc 用于 Swagger 展示班级平均成绩结构
type ClassAvgDoc struct {
	ClassName string  `json:"class_name" example:"高一1班"`
	AvgScore  float64 `json:"avg_score" example:"85.2"`
}
