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
	User         models.User           `json:"user"`
	Student      models.Student        `json:"student"`
	Course       models.Course         `json:"course"`
	Score        models.Score          `json:"score"`
	Permission   models.Permission     `json:"permission"`
	Role         models.Role           `json:"role"`
	RolePerm     models.RolePermission `json:"role_permission"`
	UserRole     models.UserRole       `json:"user_role"`
	ChartData    models.ChartData      `json:"chart_data"`
	LoginRequest models.LoginRequest   `json:"login_request"`
}
