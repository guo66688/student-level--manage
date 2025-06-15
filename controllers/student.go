// contollers/student.go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"
	"time"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

// AddStudent godoc
// @Summary 添加学生
// @Tags students
// @Accept json
// @Produce json
// @Param data body models.Student true "学生信息" example({"name": "张三", "age": 20, "class_name": "软件一班"})
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
// @Param page query int false "页码" default(1) example(1)
// @Param size query int false "每页条数" default(10) example(10)
// @Param query query string false "搜索关键词" example("张三")
// @Param student_no query string false "学号" example("20230001")
// @Param class_name query string false "班级" example("高一1班")
// @Success 200 {object} map[string]interface{} "包含学生列表和总数"
// @Failure 400 {object} map[string]string "分页参数错误"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /api/students [get]
// GetStudents godoc
// @Summary 获取学生列表
// @Tags students
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1) example(1)
// @Param size query int false "每页条数" default(10) example(10)
// @Param query query string false "搜索关键词" example("张三")
// @Param student_no query string false "学号" example("20230001")
// @Param class_name query string false "班级" example("高一1班")
// @Success 200 {object} map[string]interface{} "包含学生列表和总数"
// @Failure 400 {object} map[string]string "分页参数错误"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /api/students [get]
func GetStudents(c *gin.Context) {
	var students []models.Student
	var total int64

	// 获取分页参数
	page := c.DefaultQuery("page", "1")  // 获取页码，默认是第1页
	size := c.DefaultQuery("size", "10") // 获取每页的条数，默认10
	query := c.DefaultQuery("query", "") // 获取查询参数，默认空字符串
	studentNo := c.DefaultQuery("student_no", "") // 获取学号查询参数
	className := c.DefaultQuery("class_name", "") // 获取班级查询参数

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

	// 获取学生总数（支持查询条件）
	queryCondition := "%" + query + "%" // 用于模糊查询
	// 构造查询条件
	var conditions []string
	var args []interface{}

	// 添加查询条件
	if query != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, queryCondition)
	}
	if studentNo != "" {
		conditions = append(conditions, "student_no LIKE ?")
		args = append(args, "%"+studentNo+"%")
	}
	if className != "" {
		conditions = append(conditions, "classes.name LIKE ?")  // 使用 Class.name 作为查询条件
		args = append(args, "%" + className + "%")
	}


	// 构建查询条件
	var querySQL string
	if len(conditions) > 0 {
		querySQL = " " + strings.Join(conditions, " AND ")
	} 

	// 获取学生总数
	if err := config.DB.Model(&models.Student{}).
		Joins("LEFT JOIN classes ON classes.id = students.class_id"). // 使用 LEFT JOIN 来连接 students 和 classes
		Where(querySQL, args...).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	// 查询学生数据（分页和过滤），并加载班级名称
    if err := config.DB.Model(&models.Student{}).
		Preload("Class"). // 使用 Preload 加载关联的 Class 信息
		Joins("LEFT JOIN classes ON classes.id = students.class_id"). // 使用 JOIN 来连接
		Where(querySQL, args...). // 查询条件
		Offset((pageInt - 1) * sizeInt).Limit(sizeInt). // 分页
		Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	

	// 返回学生数据和总数
	c.JSON(http.StatusOK, gin.H{
		"data":  students,
		"total": total,
	})
}




// UpdateStudent godoc
// @Summary 更新学生信息
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "学生ID" example(1)
// @Param data body models.Student true "学生信息" example({"name": "李四", "age": 21, "class_name": "软件二班"})
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /students/{id} [put]
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
// @Param id path int true "学生ID" example(1)
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
// @Param name query string true "班级名称" example("软件一班")
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
