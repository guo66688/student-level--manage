// controllers/user.go
package controllers

import (
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// ListUsers godoc
// @Summary 获取用户列表
// @Tags user
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /users [get]
// 用户列表（支持分页与关键词）
func ListUsers(c *gin.Context) {
	var users []models.User
	var total int64

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	keyword := c.DefaultQuery("keyword", "")

	pageNum := config.StrToInt(page)
	sizeNum := config.StrToInt(size)

	offset := (pageNum - 1) * sizeNum
	db := config.DB.Model(&models.User{})
	if keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	db.Count(&total)
	db.Offset(offset).Limit(sizeNum).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"users": users,
	})
}

// AddUser godoc
// @Summary 添加用户
// @Tags user
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /users [post]
// 添加用户
func AddUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "用户名和密码不能为空"})
		return
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "user": user})
}

// UpdateUser godoc
// @Summary 更新用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /users/{id} [put]
// 更新用户基本信息
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "用户不存在"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "更新成功", "user": user})
}

// DeleteUser godoc
// @Summary 删除用户
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /users/{id} [delete]
// 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

// ResetPassword godoc
// @Summary 重置用户密码
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /users/{id}/reset_password [put]
// 重置密码
func ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "新密码不能为空"})
		return
	}
	if err := config.DB.Model(&models.User{}).Where("id = ?", id).Update("password", req.NewPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "密码重置失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "密码已重置"})
}

// GetUserPermissions godoc
// @Summary 获取当前用户权限
// @Tags user
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /user/permissions [get]
func GetUserPermissions(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
		return
	}
	user := userAny.(models.User)

	// 获取用户角色 ID
	var roleIDs []uint
	err := config.DB.
		Table("user_roles").
		Select("role_id").
		Where("user_id = ?", user.ID).
		Scan(&roleIDs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "角色查询失败"})
		return
	}

	// 查询这些角色对应的权限 key
	var keys []string
	err = config.DB.
		Table("role_permissions").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Pluck("permissions.key", &keys).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "权限查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": keys})
}

// GetUserRoles godoc
// @Summary 获取当前用户角色
// @Tags user
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /user/roles [get]
func GetUserRoles(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
		return
	}
	user := userAny.(models.User)

	var roles []models.Role
	err := config.DB.
		Joins("JOIN user_roles ur ON ur.role_id = roles.id").
		Where("ur.user_id = ?", user.ID).
		Find(&roles).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, roles)
}
