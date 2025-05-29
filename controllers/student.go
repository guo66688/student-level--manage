// contollers/student.go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"
	"time"

	"github.com/gin-gonic/gin"
)

// AddStudent godoc
// @Summary 添加学生
// @Tags students
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/students [post]
func AddStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	if err := config.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "student": student})
}

// GetStudents godoc
// @Summary 获取学生列表
// @Tags students
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/students [get]
func GetStudents(c *gin.Context) {
	var students []models.Student
	if err := config.DB.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, students)
}

// UpdateStudent godoc
// @Summary 更新学生信息
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "学生ID"
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/students/{id} [put]
func UpdateStudent(c *gin.Context) {
	var student models.Student
	id := c.Param("id")
	if err := config.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "未找到学生"})
		return
	}
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	config.DB.Save(&student)
	c.JSON(http.StatusOK, gin.H{"msg": "更新成功", "student": student})
}

// DeleteStudent godoc
// @Summary 删除学生
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "学生ID"
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/students/{id} [delete]
func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Student{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

// GetStudentsByClass godoc
// @Summary 根据班级获取学生
// @Tags students
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/students/by_class [get]
func GetStudentsByClass(c *gin.Context) {
	className := c.Query("name")
	if className == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "缺少班级名参数"})
		return
	}

	redisKey := fmt.Sprintf("class:%s:students", className)

	// 1. 先查 Redis 缓存
	if data, err := config.Redis.Get(config.Ctx, redisKey).Result(); err == nil {
		var cached []models.Student
		if json.Unmarshal([]byte(data), &cached) == nil {
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	// 2. 查询数据库
	var students []models.Student
	if err := config.DB.Where("class_name = ?", className).Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	// 3. 写入缓存
	jsonData, _ := json.Marshal(students)
	config.Redis.Set(config.Ctx, redisKey, jsonData, time.Hour)

	c.JSON(http.StatusOK, students)
}

// func invalidateClassStudentsCache(className string) {
// 	key := fmt.Sprintf("class:%s:students", className)
// 	config.Redis.Del(config.Ctx, key)
// }
