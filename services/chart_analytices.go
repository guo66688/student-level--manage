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

// UpdatePassRateChart 更新课程通过率图表数据
func UpdatePassRateChart(ctx context.Context) {
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

	// 按照课程 ID 查询每门课的通过率
	for _, cid := range courseIDs {
		var total, passed int64
		config.DB.Model(&models.Score{}).Where("course_id = ?", cid).Count(&total)
		config.DB.Model(&models.Score{}).Where("course_id = ? AND score >= 60", cid).Count(&passed)

		// 计算通过率
		rate := "0%"
		if total > 0 {
			rate = fmt.Sprintf("%.2f%%", float64(passed)/float64(total)*100)
		}

		// 将通过率数据添加到 items 列表
		items = append(items, Item{
			CourseID: cid,
			Passed:   int(passed),
			Total:    int(total),
			Rate:     rate,
		})
	}

	// 将 items 转换为符合 ChartValues 的格式
	var xValues []string
	var yValues []float64 // 将 yValues 定义为 []float64 类型
	for _, item := range items {
		// 将通过率数据转为图表格式
		xValues = append(xValues, fmt.Sprintf("课程 %d", item.CourseID)) // 使用课程 ID 作为 x 值
		yValues = append(yValues, float64(item.Passed))                // 将通过的学生数量转换为 float64，并作为 y 值
	}

	// 创建一个符合 ChartValues 的结构
	// chartValues = models.ChartValues{
	// 	X: xValues,
	// 	Y: yValues,
	// }

	// 组装图表数据
	chartData := models.ChartData{
		Type:   "pass_rate", // 图表类型：通过率
		Meta:   []models.Meta{{Key: "title", Value: "课程通过率"}},
		Values: []models.ChartValues{}, // 使用 models.ChartValues 类型
	}

	// 更新图表数据
	if err := SaveOrUpdateChartData(ctx, chartData); err != nil {
		log.Println("❌ 更新通过率图表失败:", err)
	}
}

// SaveOrUpdateChartData 保存或更新图表数据（按 type 去重）
func SaveOrUpdateChartData(ctx context.Context, data models.ChartData) error {
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
		return err
	}
	return nil
}
