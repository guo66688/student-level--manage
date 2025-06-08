// contollers/analytics.go
package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"student-level-manage/config"
	"student-level-manage/models"
	"student-level-manage/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CourseStats struct {
	CourseID uint    `json:"course_id"`
	AvgScore float64 `json:"avg_score"`
	MaxScore float64 `json:"max_score"`
	MinScore float64 `json:"min_score"`
}

// GetCourseStats godoc
// @Summary 获取课程统计
// @Tags analytics
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/analytics/course-stats [get]
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

// GetMonthlyStats godoc
// @Summary 获取月度平均成绩趋势
// @Tags analytics
// @Produce json
// @Router /api/analysis/monthly [get]
func GetMonthlyStats(c *gin.Context) {
	type Result struct {
		Month string  `json:"month"`
		Avg   float64 `json:"avg"`
	}

	var results []Result
	err := config.DB.
		Raw(`
			SELECT
			  DATE_FORMAT(exam_date, '%Y-%m') AS month,
			  AVG(score) AS avg
			FROM scores
			WHERE deleted_at IS NULL
			GROUP BY month
			ORDER BY month
		`).Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "统计失败"})
		return
	}
	c.JSON(http.StatusOK, results)
}

// GetPassRate godoc
// @Summary 获取课程通过率分析
// @Tags analytics
// @Accept json
// @Produce json
// @Param course_id query int false "课程 ID，不传则返回所有课程"
// @Success 200 {array} map[string]interface{} "返回信息"
// @Router /api/analysis/pass_rate [get]
func GetPassRate(c *gin.Context) {
	// 如果传了 course_id，则只计算单个课程
	if courseID := c.Query("course_id"); courseID != "" {
		var total, passed int64
		config.DB.Model(&models.Score{}).
			Where("course_id = ?", courseID).
			Count(&total)
		config.DB.Model(&models.Score{}).
			Where("course_id = ? AND score >= 60", courseID).
			Count(&passed)

		rate := 0.0
		if total > 0 {
			rate = float64(passed) / float64(total) * 100
		}

		c.JSON(http.StatusOK, gin.H{
			"course_id": courseID,
			"passed":    passed,
			"total":     total,
			"rate":      fmt.Sprintf("%.2f%%", rate),
		})
		return
	}

	// 不带参数时，统计所有课程
	type Pass struct {
		CourseID uint    `json:"course_id"`
		Passed   int64   `json:"passed"`
		Total    int64   `json:"total"`
		Rate     float64 `json:"rate"`
	}
	var stats []Pass

	// 先拿到所有课程 ID
	var courseIDs []uint
	config.DB.Model(&models.Course{}).Pluck("id", &courseIDs)

	// 对每个课程计算
	for _, cid := range courseIDs {
		var total, passed int64
		config.DB.Model(&models.Score{}).
			Where("course_id = ?", cid).
			Count(&total)
		config.DB.Model(&models.Score{}).
			Where("course_id = ? AND score >= 60", cid).
			Count(&passed)

		rate := 0.0
		if total > 0 {
			rate = float64(passed) / float64(total) * 100
		}

		stats = append(stats, Pass{
			CourseID: cid,
			Passed:   passed,
			Total:    total,
			Rate:     rate,
		})
	}

	c.JSON(http.StatusOK, stats)
}

// GetScoreRanking godoc
// @Summary 获取成绩排行榜（含姓名）
// @Tags analytics
// @Produce json
// @Router /api/analysis/rank [get]
func GetScoreRanking(c *gin.Context) {
	// 分页参数
	start, limit := 0, 10
	if s := c.Query("offset"); s != "" {
		fmt.Sscanf(s, "%d", &start)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	end := start + limit - 1

	// 直接连表查姓名 + 平均分
	type Row struct {
		StudentName string  `json:"student_name"`
		AvgScore    float64 `json:"avg_score"`
	}
	var all []Row
	err := config.DB.
		Table("scores AS sc").
		Select("st.name AS student_name, AVG(sc.score) AS avg_score").
		Joins("JOIN students AS st ON sc.student_id = st.id").
		Group("sc.student_id").
		Order("avg_score DESC").
		Scan(&all).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "数据库查询失败"})
		return
	}

	// 分页返回
	if start > len(all)-1 {
		c.JSON(http.StatusOK, []Row{})
		return
	}
	if end >= len(all) {
		end = len(all) - 1
	}
	c.JSON(http.StatusOK, all[start:end+1])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ClearScoreRankingCache godoc
// @Summary 清空成绩排行榜缓存
// @Tags analytics
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/analysis/rank/cache [delete]
func ClearScoreRankingCache(c *gin.Context) {
	err := config.Redis.Del(config.Ctx, "score:rank").Err()
	if err != nil {
		c.JSON(500, gin.H{"msg": "清除失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "缓存已清除"})
}

// GetChartFromMongo godoc
// @Summary 获取图表数据
// @Tags charts
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/charts/data [get]
func GetChartFromMongo(c *gin.Context) {
	chartType := c.Query("type")
	if chartType == "" {
		c.JSON(400, gin.H{"msg": "缺少参数 type"})
		return
	}
	data, err := services.GetChartData(chartType)
	if err != nil {
		c.JSON(500, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(200, data)
}

// GetChartTypes godoc
// @Summary 图表类型列表
// @Tags charts
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/charts/types [get]
func GetChartTypes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("analysis").Collection("charts")

	// 只查询 type 字段的去重值
	cursor, err := collection.Distinct(ctx, "type", bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	var types []string
	for _, v := range cursor {
		if str, ok := v.(string); ok {
			types = append(types, str)
		}
	}

	c.JSON(http.StatusOK, types)
}

// DeleteChartData godoc
// @Summary 删除图表数据
// @Tags charts
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/charts/data [delete]
func DeleteChartData(c *gin.Context) {
	chartType := c.Query("type")
	if chartType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "缺少 type 参数"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("analysis").Collection("charts")
	res, err := collection.DeleteMany(ctx, bson.M{"type": chartType})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "删除成功",
		"deleted": res.DeletedCount,
	})
}

// GetDashboardStats godoc
// @Summary 仪表盘统计数据
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/dashboard/stats [get]
func GetDashboardStats(c *gin.Context) {
	var studentsCount int64
	var coursesCount int64
	var classesCount int64
	var avgScore float64

	config.DB.Model(&models.Student{}).Count(&studentsCount)
	config.DB.Model(&models.Course{}).Count(&coursesCount)
	config.DB.Model(&models.Class{}).Count(&classesCount)
	// 这里假设 class 是你存班级的表
	config.DB.
		Table("scores").
		Select("AVG(score)").
		Row().
		Scan(&avgScore)

	c.JSON(http.StatusOK, gin.H{
		"students": studentsCount,
		"courses":  coursesCount,
		"classes":  classesCount,
		"avgScore": avgScore,
	})
}
