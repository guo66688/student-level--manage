// services/static_chart_data_generator.go
package services

import (
	"context"
	"fmt"
	"student-level-manage/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// StaticChartDataGenerator 用于生成静态图表的数据
type StaticChartDataGenerator struct {
	// 数据库连接字段
	db *mongo.Database
}

// NewStaticChartDataGenerator 构造函数，用于初始化 StaticChartDataGenerator
func NewStaticChartDataGenerator(db *mongo.Database) *StaticChartDataGenerator {
	return &StaticChartDataGenerator{db: db}
}

// 实现 ChartDataGenerator 接口中的 GenerateData 方法
func (g *StaticChartDataGenerator) GenerateData() ([]map[string]interface{}, error) {
	// 由于静态数据应该已经存在于数据库中
	return nil, fmt.Errorf("GenerateData 方法未实现，静态数据需根据 ID 获取")
}

// 根据ID获取数据的方法，更新返回类型为 []map[string]interface{}
func (g *StaticChartDataGenerator) GenerateDataByID(chartID string) ([]map[string]interface{}, error) {
	// 获取图表数据
	chartData, err := g.getChartDataByID(chartID)
	if err != nil {
		fmt.Println("Error fetching chart data:", err)
		return nil, err
	}

	// 打印获取到的 chartData，检查它的内容
	fmt.Println("Found chart data:", chartData)

	var xValues []string
	var yValues []float64

	// 遍历 chartData.Values 来提取 x 和 y 数据
	for _, value := range chartData.Values {
		// 获取每个 ChartValues 对象的 X 和 Y
		x := value.X
		y := value.Y

		// 打印获取到的 x 和 y 数组
		fmt.Println("x values:", x)
		fmt.Println("y values:", y)

		if len(x) == len(y) {
			// 遍历 x 和 y 数组
			for i := 0; i < len(x); i++ {
				// 确保 x 是字符串类型，y 是数字类型
				xValues = append(xValues, x[i]) // 将 x 值添加到 xValues 列表
				fmt.Printf("Added x value: %s\n", x[i])

				yValues = append(yValues, y[i]) // 将 y 值添加到 yValues 列表
				fmt.Printf("Added y value: %f\n", y[i])
			}
		} else {
			fmt.Println("Mismatch between x and y array lengths")
			return nil, fmt.Errorf("x 和 y 数组长度不匹配")
		}
	}

	// 打印最终的解析结果
	fmt.Println("Final parsed data - x values:", xValues, "y values:", yValues)

	// 返回解析后的 x 和 y 值，满足接口的返回要求
	return []map[string]interface{}{
		{
			"x": xValues,
			"y": yValues,
		},
	}, nil
}

// getChartDataByID 根据 chartID 从数据库获取图表数据
func (g *StaticChartDataGenerator) getChartDataByID(chartID string) (*models.ChartData, error) {
	// 构造查询条件
	objectID, err := primitive.ObjectIDFromHex(chartID)
	if err != nil {
		return nil, fmt.Errorf("无效的 chartID: %v", err)
	}

	// 打印转换后的 ObjectID，检查是否正确
	fmt.Println("Searching for chart data with ObjectID:", objectID)

	// 获取 collection
	collection := g.db.Collection("charts") // 假设你有一个 "charts" 集合

	// 查询数据
	var chartData models.ChartData
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&chartData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 打印没有找到记录的错误，帮助调试
			fmt.Println("No documents found for chartID:", chartID)
			return nil, fmt.Errorf("未找到图表数据")
		}
		return nil, fmt.Errorf("查询图表数据失败: %v", err)
	}

	// 打印查询结果，确认返回数据
	fmt.Println("Found chart data:", chartData)

	return &chartData, nil
}
