// controllers/permission.go
package controllers

import (
	"fmt"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// GetPermissions godoc
// @Summary 获取权限列表
// @Tags permission
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/permissions [get]
// 获取所有权限
func GetPermissions(c *gin.Context) {
	var list []models.Permission
	if err := config.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// AddPermission godoc
// @Summary 添加权限
// @Tags permission
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/permissions [post]
// 添加权限
func AddPermission(c *gin.Context) {
	var perm models.Permission
	if err := c.ShouldBindJSON(&perm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	if err := config.DB.Create(&perm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "id": perm.ID})
}

// DeletePermission godoc
// @Summary 删除权限
// @Tags permission
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/permissions/{id} [delete]
// 删除权限
func DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Permission{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

// SetRolePermissions godoc
// @Summary 设置角色权限
// @Tags permission
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/role/permissions [post]
// 设置角色权限
func SetRolePermissions(c *gin.Context) {
	var req struct {
		RoleID        uint   `json:"role_id"`
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}

	// 先删除旧的关联
	if err := config.DB.Where("role_id = ?", req.RoleID).Delete(&models.RolePermission{}).Error; err != nil {
		c.JSON(500, gin.H{"msg": "删除旧权限失败"})
		return
	}

	// 添加新关联
	var records []models.RolePermission
	for _, pid := range req.PermissionIDs {
		records = append(records, models.RolePermission{
			RoleID:       req.RoleID,
			PermissionID: pid,
		})
	}
	if len(records) > 0 {
		if err := config.DB.Create(&records).Error; err != nil {
			c.JSON(500, gin.H{"msg": "添加权限失败"})
			return
		}
	}

	// 清除所有绑定该角色的用户的权限缓存
	var userIDs []uint
	config.DB.
		Table("user_roles").
		Select("user_id").
		Where("role_id = ?", req.RoleID).
		Scan(&userIDs)

	for _, uid := range userIDs {
		config.Redis.Del(config.Ctx, fmt.Sprintf("user:%d:permissions", uid))
	}

	c.JSON(200, gin.H{"msg": "角色权限已更新"})
}
