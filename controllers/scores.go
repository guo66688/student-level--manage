package controllers

import (
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

func AddScore(c *gin.Context) {
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
}

func GetScores(c *gin.Context) {
	var scores []models.Score
	if err := config.DB.Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, scores)
}
