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
	"github.com/redis/go-redis/v9"
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
// @Summary 获取月度成绩统计
// @Tags analytics
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/analysis/monthly [get]
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

// GetPassRate godoc
// @Summary 获取课程通过率分析
// @Tags analytics
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/analysis/pass_rate [get]
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

// GetScoreRanking godoc
// @Summary 获取成绩排行榜
// @Tags analytics
// @Accept json
// @Produce json
// @Param data body object true "请求参数"
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /api/analysis/rank [get]
// 获取学生成绩排行榜（先查 Redis，缓存失效则查 DB）
func GetScoreRanking(c *gin.Context) {
	key := "score:rank"

	// 分页参数
	start := 0
	limit := 10
	if s := c.Query("offset"); s != "" {
		fmt.Sscanf(s, "%d", &start)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	end := start + limit - 1

	// 从 Redis 获取
	zset, err := config.Redis.ZRevRangeWithScores(config.Ctx, key, int64(start), int64(end)).Result()
	if err == nil && len(zset) > 0 {
		var result []gin.H
		for _, z := range zset {
			result = append(result, gin.H{
				"student_id": z.Member,
				"avg_score":  z.Score,
			})
		}
		c.JSON(200, result)
		return
	}

	// 数据库查询
	type Result struct {
		StudentID uint    `json:"student_id"`
		AvgScore  float64 `json:"avg_score"`
	}
	var list []Result
	if err := config.DB.
		Table("scores").
		Select("student_id, AVG(score) as avg_score").
		Group("student_id").
		Order("avg_score DESC").
		Scan(&list).Error; err != nil {
		c.JSON(500, gin.H{"msg": "数据库查询失败"})
		return
	}

	// 写入 Redis ZSet 缓存
	pipe := config.Redis.Pipeline()
	for _, s := range list {
		pipe.ZAdd(config.Ctx, key, redis.Z{
			Score:  s.AvgScore,
			Member: s.StudentID,
		})
	}
	pipe.Expire(config.Ctx, key, 3600) // 1 小时
	_, _ = pipe.Exec(config.Ctx)

	// 返回分页部分
	end = min(end+1, len(list))
	c.JSON(200, list[start:end])
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
