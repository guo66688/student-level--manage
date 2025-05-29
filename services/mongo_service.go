package services

import (
	"context"
	"student-level-manage/config"
	"student-level-manage/models"
)

func SaveChartData(data models.ChartData) error {
	coll := config.MongoClient.Database("student_db").Collection("charts")
	_, err := coll.InsertOne(context.TODO(), data)
	return err
}
