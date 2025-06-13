// contollers/courses.go
package controllers

import (
	"net/http"
	"strconv"

	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// AddCourse godoc
// @Summary 添加课程
// @Tags courses
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/courses [post]
func AddCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	// ⚠️ 先校验
	if course.CourseName == "" || course.CourseCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "课程名和课程代码不能为空"})
		return
	}

	// ✅ 再插入数据库
	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "course": course})
}

// GetCourses godoc
// @Summary 获取课程列表
// @Tags courses
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/courses [get]
func GetCourses(c *gin.Context) {
	var courses []models.Course
	var total int64

	// 获取分页参数
	page := c.DefaultQuery("page", "1")  // 获取页码，默认是第1页
	size := c.DefaultQuery("size", "10") // 获取每页的条数，默认10

	// 将字符串转换为整数
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "页码无效"})
		return
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "每页条数无效"})
		return
	}

	// 获取课程总数
	if err := config.DB.Model(&models.Course{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	// 查询课程数据（分页）
	if err := config.DB.Offset((pageInt - 1) * sizeInt).Limit(sizeInt).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	// 返回课程数据和总数
	c.JSON(http.StatusOK, gin.H{
		"data":  courses,
		"total": total,
	})
}

// @Summary 更新课程
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "课程ID"
// @Param course body models.Course true "课程信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/courses/{id} [put]
func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course
	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "未找到课程"})
		return
	}
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	config.DB.Save(&course)
	c.JSON(http.StatusOK, gin.H{"msg": "更新成功", "course": course})
}

// @Summary 删除课程
// @Tags courses
// @Param id path int true "课程ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/courses/{id} [delete]
func DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Course{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
