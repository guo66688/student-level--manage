// models\doc_analytics.go
package models

// CourseStats 用于课程统计响应
type CourseStats struct {
	CourseID uint    `json:"course_id" example:"101"`
	AvgScore float64 `json:"avg_score" example:"85.3"`
	MaxScore float64 `json:"max_score" example:"98"`
	MinScore float64 `json:"min_score" example:"62"`
}

// Result 用于月度平均成绩
type Result struct {
	Month string  `json:"month" example:"2024-06"`
	Avg   float64 `json:"avg" example:"82.7"`
}

// Pass 用于课程通过率
type Pass struct {
	CourseID uint    `json:"course_id" example:"101"`
	Passed   int64   `json:"passed" example:"45"`
	Total    int64   `json:"total" example:"50"`
	Rate     float64 `json:"rate" example:"90"`
}

// ScoreRankRow 用于成绩排行榜响应
type ScoreRankRow struct {
	StudentName string  `json:"student_name" example:"张三"`
	AvgScore    float64 `json:"avg_score" example:"87.5"`
}
