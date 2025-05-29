package controllers

import (
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

func AddClass(c *gin.Context) {
	var class models.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	if err := config.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "class": class})
}

func GetClassStudents(c *gin.Context) {
	className := c.Param("name")
	cacheKey := "class:" + className + ":students"

	// 先查 Redis 缓存
	if data, err := config.Redis.Get(config.Ctx, cacheKey).Result(); err == nil {
		c.JSON(http.StatusOK, gin.H{"from": "redis", "students": data})
		return
	}

	// 查 MySQL
	var students []models.Student
	if err := config.DB.Where("class_name = ?", className).Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	// 存入 Redis
	_ = config.Redis.Set(config.Ctx, cacheKey, students, 0).Err()

	c.JSON(http.StatusOK, gin.H{"from": "mysql", "students": students})
}
