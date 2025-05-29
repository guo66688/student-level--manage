package services

import (
	"context"
	"log"
	"student-level-manage/config"
	"student-level-manage/models"
	"time"
)

func SaveChartData(data models.ChartData) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.MongoDB.Database("analysis").Collection("charts").InsertOne(ctx, data)
	if err != nil {
		log.Println("❌ MongoDB 图表数据保存失败:", err)
	}
}
