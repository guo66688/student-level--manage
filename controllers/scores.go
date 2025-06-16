// contollers/scores.go
package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"student-level-manage/config"
	"student-level-manage/models"
	"student-level-manage/redisop"
	"student-level-manage/services"

	"github.com/gin-gonic/gin"
)

// AddScore godoc
// @Summary 添加成绩
// @Tags scores
// @Accept json
// @Produce json
// @Param data body models.Score true "成绩信息" example({"student_id": 1, "course_id": 101, "score": 85, "exam_date": "2024-06-01"})
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /scores [post]
func AddScore(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req models.Score
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	if req.Score < 0 || req.Score > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "分数必须在 0 到 100 之间"})
		return
	}
	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加成绩失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "score": req})

	// 添加成绩后更新 Redis 排行榜
	var avg float64
	config.DB.
		Table("scores").
		Select("AVG(score)").
		Where("student_id = ?", req.StudentID).
		Scan(&avg)

	_ = redisop.UpdateStudentRank(req.StudentID, avg) // 忽略错误处理

	// ✅ 自动更新 Mongo 图表（按月统计）
	var results []struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}
	config.DB.
		Raw(`
			SELECT DATE_FORMAT(exam_date, '%Y-%m') as month, COUNT(*) as count
			FROM scores
			GROUP BY month
			ORDER BY month
		`).Scan(&results)

	services.SaveOrUpdateChartData(ctx, models.ChartData{
		Type:   "monthly",
		Meta:   map[string]interface{}{},
		Values: results,
	})

	// ✅ 自动更新课程通过率图表
	services.UpdatePassRateChart(ctx)
}

// GetScores godoc
// @Summary 获取成绩列表
// @Tags scores
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1) example(1)
// @Param size query int false "每页条数" default(10) example(10)
// @Param query query string false "搜索关键词" example("张三")
// @Param student_no query string false "学号" example("20230001")
// @Param course_name query string false "课程" example("数学")
// @Success 200 {object} map[string]interface{} "包含成绩列表和总数"
// @Failure 400 {object} map[string]string "分页参数错误"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /scores [get]
func GetScores(c *gin.Context) {
	var scores []ScoreWithDetails
	var total int64
	log.Println("Before Query Execution")

	// 获取分页参数
	page := c.DefaultQuery("page", "1")             // 获取页码，默认是第1页
	size := c.DefaultQuery("size", "10")            // 获取每页的条数，默认10
	query := c.DefaultQuery("query", "")            // 获取查询参数，默认空字符串
	studentNo := c.DefaultQuery("student_no", "")   // 获取学号查询参数
	courseName := c.DefaultQuery("course_name", "") // 获取课程查询参数

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

	// 获取成绩总数（支持查询条件）
	queryCondition := "%" + query + "%" // 用于模糊查询
	// 构造查询条件
	var conditions []string
	var args []interface{}

	// 添加查询条件
	if query != "" {
		conditions = append(conditions, "students.name LIKE ?")
		args = append(args, queryCondition)
	}
	if studentNo != "" {
		conditions = append(conditions, "students.student_no LIKE ?")
		args = append(args, "%"+studentNo+"%")
	}
	if courseName != "" {
		conditions = append(conditions, "courses.course_name LIKE ?") // 使用 Course.course_name 作为查询条件
		args = append(args, "%"+courseName+"%")
	}

	// 构建查询条件
	var querySQL string
	if len(conditions) > 0 {
		querySQL = " " + strings.Join(conditions, " AND ")
	}

	// 获取成绩总数
	if err := config.DB.Model(&models.Score{}).
		Joins("LEFT JOIN students ON students.id = scores.student_id").
		Joins("LEFT JOIN courses ON courses.id = scores.course_id").
		Where(querySQL, args...).
		Count(&total).Error; err != nil {
		fmt.Println("Error during count query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	fmt.Println("Total Scores:", total)

	fmt.Println("Offset:", (pageInt-1)*sizeInt)
	fmt.Println("Limit:", sizeInt)

	// 查询成绩数据（分页和过滤），并确保学生姓名和课程名称能够正确加载
	err = config.DB.Model(&models.Score{}).
		Joins("LEFT JOIN students ON students.id = scores.student_id").
		Joins("LEFT JOIN courses ON courses.id = scores.course_id").
		Where(querySQL, args...).
		Select("scores.id, scores.student_id, scores.course_id, scores.exam_date, scores.score, students.name AS student_name, students.student_no AS student_no, courses.course_name AS course_name, courses.course_code AS course_code, scores.created_at, scores.updated_at").
		Offset((pageInt - 1) * sizeInt).Limit(sizeInt). // 确保分页参数被正确传递
		Scan(&scores).Error                             // 使用 Scan 显式映射查询结果到结构体
	if err != nil {
		fmt.Println("Error during data query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	fmt.Println("Scores Data:", scores)
	fmt.Println("Scores Data:", total)
	// 返回包含成绩、学生和课程数据的完整列表和总数
	c.JSON(http.StatusOK, gin.H{
		"data":  scores,
		"total": total,
	})
}

type ScoreWithDetails struct {
	ID          uint    `json:"id"`
	StudentID   uint    `json:"student_id"`
	CourseID    uint    `json:"course_id"`
	ExamDate    string  `json:"exam_date"`
	Score       float64 `json:"score"`
	StudentName string  `json:"student_name"`
	StudentNo   string  `json:"student_no"`
	CourseName  string  `json:"course_name"`
	CourseCode  string  `json:"course_code"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// UpdateScore godoc
// @Summary 更新成绩
// @Tags scores
// @Accept json
// @Produce json
// @Param id path int true "成绩ID" example(1)
// @Param data body models.Score true "成绩信息" example({"student_id": 1, "course_id": 101, "score": 90, "exam_date": "2024-06-01"})
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /scores/{id} [put]
func UpdateScore(c *gin.Context) {
	id := c.Param("id")
	var score models.Score
	if err := config.DB.First(&score, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "未找到成绩"})
		return
	}
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	config.DB.Save(&score)
	c.JSON(http.StatusOK, gin.H{"msg": "更新成功", "score": score})
}

// DeleteScore godoc
// @Summary 删除成绩
// @Tags scores
// @Accept json
// @Produce json
// @Param id path int true "成绩ID" example(1)
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /scores/{id} [delete]
func DeleteScore(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Score{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
