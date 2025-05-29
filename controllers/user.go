// controllers/user.go
package controllers

import (
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

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

// 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

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
