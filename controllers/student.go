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

func GetStudents(c *gin.Context) {
	var students []models.Student
	if err := config.DB.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, students)
}

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

func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Student{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

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
func invalidateClassStudentsCache(className string) {
	key := fmt.Sprintf("class:%s:students", className)
	config.Redis.Del(config.Ctx, key)
}
