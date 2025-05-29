package services

import (
	"context"
	"student-level-manage/config"
	"student-level-manage/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetChartData(chartType string) ([]models.ChartData, error) {
	collection := config.MongoDB.Collection("charts")
	cur, err := collection.Find(context.Background(), bson.M{"type": chartType})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var results []models.ChartData
	for cur.Next(context.Background()) {
		var item models.ChartData
		if err := cur.Decode(&item); err == nil {
			results = append(results, item)
		}
	}
	return results, nil
}
