package controllers

import (
	"net/http"
	"student-level-manage/config"

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
