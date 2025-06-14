// models\doc_models.go
package models

// CourseDoc 用于 Swagger 文档展示
type CourseDoc struct {
	ID         uint    `json:"id" example:"1"`
	CreatedAt  string  `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt  string  `json:"updated_at" example:"2024-01-01T12:00:00Z"`
	CourseName string  `json:"course_name" example:"数学"`
	CourseCode string  `json:"course_code" example:"MATH101"`
	Credit     float64 `json:"credit" example:"3"`
	TeacherID  uint    `json:"teacher_id" example:"10"`
}

// UserDoc 用于 Swagger 文档展示
type UserDoc struct {
	ID        uint   `json:"id" example:"1"`
	CreatedAt string `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T12:00:00Z"`
	Username  string `json:"username" example:"admin"`
	Password  string `json:"password" example:"123456"`
	Nickname  string `json:"nickname" example:"管理员"`
	IsActive  bool   `json:"is_active" example:"true"`
}

// StudentDoc 用于 Swagger 文档展示
type StudentDoc struct {
	ID        uint   `json:"id" example:"1"`
	CreatedAt string `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T12:00:00Z"`
	StudentNo string `json:"student_no" example:"20230001"`
	Name      string `json:"name" example:"张三"`
	Gender    string `json:"gender" example:"男"`
	BirthDate string `json:"birth_date" example:"2005-09-01"`
	ClassID   uint   `json:"class_id" example:"2"`
}

// ScoreDoc 用于 Swagger 文档展示
type ScoreDoc struct {
	ID        uint    `json:"id" example:"1"`
	CreatedAt string  `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2024-01-01T12:00:00Z"`
	StudentID uint    `json:"student_id" example:"1"`
	CourseID  uint    `json:"course_id" example:"101"`
	ExamDate  string  `json:"exam_date" example:"2025-05-01"`
	Score     float64 `json:"score" example:"89.5"`
}

// ClassDoc 用于 Swagger 文档展示
type ClassDoc struct {
	ID        uint   `json:"id" example:"1"`
	CreatedAt string `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T12:00:00Z"`
	Name      string `json:"name" example:"高一1班"`
}

type RoleDoc struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"管理员"`
	CreatedAt string `json:"created_at" example:"2024-01-01T12:00:00Z"`
}

type PermissionDoc struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"学生管理"`
	Key  string `json:"key" example:"student:view"`
}
