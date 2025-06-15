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

// GetCourseStats godoc
// @Summary 获取课程统计
// @Description 获取每门课程的平均、最高和最低成绩
// @Tags analytics
// @Produce json
// @Success 200 {array} models.CourseStats "课程统计列表"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /analytics/course-stats [get]
func GetCourseStats(c *gin.Context) {
	var stats []models.CourseStats
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
// @Description 获取每月平均成绩的时间序列数据
// @Tags analytics
// @Produce json
// @Success 200 {array} models.Result "月度平均成绩列表"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /analysis/monthly [get]
func GetMonthlyStats(c *gin.Context) {

	var results []models.Result
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
// @Description 可选传入 course_id 获取单个课程，否则返回所有课程通过率
// @Tags analytics
// @Produce json
// @Param course_id query int false "课程 ID（可选，单个课程）" example(1001)
// @Success 200 {array} models.Pass "通过率列表或单个课程分析"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /analysis/pass_rate [get]
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

	var stats []models.Pass

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

		stats = append(stats, models.Pass{
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
// @Description 获取学生按平均成绩排名的前10名
// @Tags analytics
// @Produce json
// @Param offset query int false "起始位置（默认0）" example(0)
// @Param limit query int false "返回数量（默认10）" example(10)
// @Success 200 {array} models.ScoreRankRow "学生平均成绩排行榜"
// @Failure 500 {object} map[string]string "数据库查询失败"
// @Router /analysis/rank [get]
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

	var all []models.ScoreRankRow
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
		c.JSON(http.StatusOK, []models.ScoreRankRow{})
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
// @Success 200 {object} map[string]string "清除结果"
// @Failure 500 {object} map[string]string "清除失败"
// @Router /analysis/rank/cache [delete]
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
// @Produce json
// @Param type query string true "图表类型" example("score-trend")
// @Success 200 {object} object "图表数据内容"
// @Failure 400 {object} map[string]string "缺少参数"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /charts/data [get]
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
// @Produce json
// @Success 200 {array} string "图表类型集合"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /charts/types [get]
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
// @Produce json
// @Param type query string true "图表类型" example("pass-rate")
// @Success 200 {object} map[string]interface{} "删除结果"
// @Failure 400 {object} map[string]string "参数缺失"
// @Failure 500 {object} map[string]string "删除失败"
// @Router /charts/data [delete]
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
// @Description 包括学生数、课程数、班级数、平均成绩
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{} "统计数据"
// @Router /dashboard/stats [get]
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


// GetExamCount godoc
// @Summary 获取考试总次数
// @Description 按 course_id 和 exam_date 去重统计考试次数（即课程+考试日视为一次考试）
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]int "格式: { \"total\": 23 }"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /analysis/exam_count [get]
func GetExamCount(c *gin.Context) {
	var count int64
	err := config.DB.
		Model(&models.Score{}).
		Select("COUNT(DISTINCT course_id, exam_date)").
		Count(&count).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(200, gin.H{"total": count})
}


// GetClassAvgScore godoc
// @Summary 获取班级平均成绩
// @Description 获取每个班级的平均成绩，用于图表展示
// @Tags analytics
// @Produce json
// @Success 200 {array} models.ClassAvgDoc "班级平均成绩数组"
// @Failure 500 {object} map[string]string "数据库查询失败"
// @Router /analysis/class_avg [get]
func GetClassAvgScore(c *gin.Context) {
	type Result struct {
		ClassName string  `json:"class_name"`
		AvgScore  float64 `json:"avg_score"`
	}

	var results []Result

	err := config.DB.
		Table("scores").
		Select("classes.name AS class_name, AVG(scores.score) AS avg_score").
		Joins("JOIN students ON scores.student_id = students.id").
		Joins("JOIN classes ON students.class_id = classes.id").
		Group("classes.name").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}
	c.JSON(http.StatusOK, results)
}


