// contollers/scores.go
package controllers

import (
	"context"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"
	"student-level-manage/redisop"
	"student-level-manage/services"
	"time"

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
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /scores [get]
func GetScores(c *gin.Context) {
	var scores []models.Score
	if err := config.DB.Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, scores)
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
