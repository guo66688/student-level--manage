package controllers

import (
	"fmt"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"
	"student-level-manage/redisop"

	"github.com/gin-gonic/gin"
)

type CourseStats struct {
	CourseID uint    `json:"course_id"`
	AvgScore float64 `json:"avg_score"`
	MaxScore float64 `json:"max_score"`
	MinScore float64 `json:"min_score"`
}

func GetCourseStats(c *gin.Context) {
	var stats []CourseStats
	query := `
		SELECT 
			course_id,
			AVG(score) AS avg_score,
			MAX(score) AS max_score,
			MIN(score) AS min_score
		FROM scores
		WHERE deleted_at IS NULL
		GROUP BY course_id
	`
	if err := config.DB.Raw(query).Scan(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// 学生成绩排名接口
func GetScoreRanking(c *gin.Context) {
	type Result struct {
		StudentID uint    `json:"student_id"`
		AvgScore  float64 `json:"avg_score"`
	}

	var results []Result
	err := config.DB.
		Table("scores").
		Select("student_id, AVG(score) as avg_score").
		Group("student_id").
		Order("avg_score DESC").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, results)
}

// 成绩按月份统计图表
func GetMonthlyStats(c *gin.Context) {
	type Result struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}

	var results []Result
	err := config.DB.
		Raw(`
			SELECT DATE_FORMAT(exam_date, '%Y-%m') as month, COUNT(*) as count
			FROM scores
			GROUP BY month
			ORDER BY month
		`).Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "统计失败"})
		return
	}
	c.JSON(http.StatusOK, results)
}

// 某科考试通过率分析

func GetPassRate(c *gin.Context) {
	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "缺少 course_id 参数"})
		return
	}

	var total, passed int64
	config.DB.Model(&models.Score{}).Where("course_id = ?", courseID).Count(&total)
	config.DB.Model(&models.Score{}).Where("course_id = ? AND score >= 60", courseID).Count(&passed)

	if total == 0 {
		c.JSON(http.StatusOK, gin.H{"msg": "暂无数据", "rate": "0%"})
		return
	}

	rate := float64(passed) / float64(total) * 100
	c.JSON(http.StatusOK, gin.H{
		"course_id": courseID,
		"passed":    passed,
		"total":     total,
		"rate":      fmt.Sprintf("%.2f%%", rate),
	})
}

func GetScoreRankingRedis(c *gin.Context) {
	ranks, err := redisop.GetTopN(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取 Redis 排行榜失败"})
		return
	}

	var result []gin.H
	for _, z := range ranks {
		result = append(result, gin.H{
			"student_id": z.Member,
			"avg_score":  z.Score,
		})
	}
	c.JSON(http.StatusOK, result)
}
