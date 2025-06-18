// services/chart.go
package services

import (
	"context"
	"fmt"
	"log"

	"student-level-manage/config"
	"student-level-manage/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// 保存或更新图表数据（按 type 去重）
func SaveOrUpdateChartData(ctx context.Context, data models.ChartData) {
	filter := bson.M{"type": data.Type}
	update := bson.M{
		"$set": bson.M{
			"type":   data.Type,
			"meta":   data.Meta,
			"values": data.Values,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := config.MongoClient.
		Database("analysis").
		Collection("charts").
		UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("❌ MongoDB 更新图表数据失败:", err)
	}
}

func UpdatePassRateChart(ctx context.Context) {
	// 不再需要自己创建 ctx
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// 查询每门课的通过率（分数 >= 60）
	type Item struct {
		CourseID uint   `json:"course_id" bson:"course_id"`
		Passed   int    `json:"passed" bson:"passed"`
		Total    int    `json:"total" bson:"total"`
		Rate     string `json:"rate" bson:"rate"`
	}
	var items []Item

	// 查询所有课程 ID
	var courseIDs []uint
	config.DB.Table("scores").Select("DISTINCT course_id").Scan(&courseIDs)

	for _, cid := range courseIDs {
		var total, passed int64
		config.DB.Model(&models.Score{}).Where("course_id = ?", cid).Count(&total)
		config.DB.Model(&models.Score{}).Where("course_id = ? AND score >= 60", cid).Count(&passed)

		rate := "0%"
		if total > 0 {
			rate = fmt.Sprintf("%.2f%%", float64(passed)/float64(total)*100)
		}
		items = append(items, Item{
			CourseID: cid,
			Passed:   int(passed),
			Total:    int(total),
			Rate:     rate,
		})
	}

	SaveOrUpdateChartData(ctx, models.ChartData{
		Type:   "pass_rate",
		Meta:   map[string]interface{}{},
		Values: items,
	})
}
