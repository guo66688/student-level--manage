// controllers/role.go
package controllers

import (
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// 获取角色列表
func GetRoles(c *gin.Context) {
	var roles []models.Role
	if err := config.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// 添加角色
func AddRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	if err := config.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功"})
}

// 删除角色
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

// 获取某个角色绑定的权限 ID 列表
func GetPermissionsByRole(c *gin.Context) {
	roleID := c.Query("role_id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "缺少 role_id 参数"})
		return
	}

	var ids []uint
	err := config.DB.
		Table("role_permissions").
		Select("permission_id").
		Where("role_id = ?", roleID).
		Scan(&ids).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, ids)
}
