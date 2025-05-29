package controllers

import (
	"fmt"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
// func GetScoreRanking(c *gin.Context) {
// 	type Result struct {
// 		StudentID uint    `json:"student_id"`
// 		AvgScore  float64 `json:"avg_score"`
// 	}

// 	var results []Result
// 	err := config.DB.
// 		Table("scores").
// 		Select("student_id, AVG(score) as avg_score").
// 		Group("student_id").
// 		Order("avg_score DESC").
// 		Scan(&results).Error

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, results)
// }

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

func ClearScoreRankingCache(c *gin.Context) {
	err := config.Redis.Del(config.Ctx, "score:rank").Err()
	if err != nil {
		c.JSON(500, gin.H{"msg": "清除失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "缓存已清除"})
}
