package controllers

import (
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// 设置用户角色
func SetUserRole(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id"`
		RoleID uint `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	// 删除旧绑定
	config.DB.Where("user_id = ?", req.UserID).Delete(&models.UserRole{})

	// 插入新绑定
	ur := models.UserRole{UserID: req.UserID, RoleID: req.RoleID}
	if err := config.DB.Create(&ur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "绑定失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "绑定成功"})
}

// 查询用户的角色 ID
func GetUserRole(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "缺少 user_id 参数"})
		return
	}

	var roleID uint
	err := config.DB.Table("user_roles").
		Select("role_id").
		Where("user_id = ?", userID).
		Scan(&roleID).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role_id": roleID})
}
