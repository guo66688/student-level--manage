// controllers\doc_controller.go
package controllers

import (
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// @Summary [文档] 展示所有模型结构体（用于 Swagger 展示）
// @Description 此接口无实际用途，仅用于强制 Swagger 输出所有模型结构体定义
// @Tags ModelDocs
// @Accept json
// @Produce json
// @Success 200 {object} docModels
// @Router /doc/models [get]
func ShowAllModels(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "模型文档展示接口，无实际返回"})
}

// docModels 用于聚合展示所有模型结构体定义
type docModels struct {
	User         models.UserDoc         `json:"user"`
	Student      models.StudentDoc      `json:"student"`
	Course       models.CourseDoc       `json:"course"`
	Score        models.ScoreDoc        `json:"score"`
	Class        models.ClassDoc        `json:"class"`
	Permission   models.Permission      `json:"permission"`
	Role         models.Role            `json:"role"`
	RolePerm     models.RolePermission  `json:"role_permission"`
	UserRole     models.UserRole        `json:"user_role"`
	ChartData    models.ChartData       `json:"chart_data"`
	LoginRequest models.LoginRequest    `json:"login_request"`
	RoleDoc       models.RoleDoc       `json:"role_doc"`
	PermissionDoc models.PermissionDoc `json:"permission_doc"`
	CourseStats    models.CourseStats    `json:"course_stats"`
	Result         models.Result         `json:"result"`
	Pass           models.Pass           `json:"pass"`
	ScoreRankRow   models.ScoreRankRow   `json:"score_rank"`
	ClassAvg models.ClassAvgDoc `json:"class_avg"`
}
