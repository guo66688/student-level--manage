// services/static_chart_data_generator.go
package services

import (
	"context"
	"fmt"

	"student-level-manage/config"
	"student-level-manage/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// StaticChartDataGenerator 用于生成静态图表的数据
type StaticChartDataGenerator struct{}

// 根据ID获取数据的方法
func (g *StaticChartDataGenerator) GenerateDataByID(id string) ([]map[string]interface{}, error) {
	// 将传入的 id 字符串转换为 ObjectID 类型
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的图表 ID 格式: %v", err)
	}

	// 假设 MongoDB 中有一个名为 "charts" 的集合存储图表数据
	collection := config.MongoClient.Database("analysis").Collection("charts")

	// 在 MongoDB 中查找图表数据
	var chartData models.ChartData
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&chartData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("未找到ID为 %s 的图表数据", id)
		}
		return nil, fmt.Errorf("查询图表数据时发生错误: %v", err)
	}

	// 假设返回的 chartData 是类似 {"x": "2024-01", "y": 80.0} 这样的数据结构
	var values []map[string]interface{}
	for _, v := range chartData.Values.([]interface{}) {
		// 根据实际存储格式对数据进行适当处理
		values = append(values, map[string]interface{}{"x": v.(map[string]interface{})["x"], "y": v.(map[string]interface{})["y"]})
	}

	return values, nil
}
