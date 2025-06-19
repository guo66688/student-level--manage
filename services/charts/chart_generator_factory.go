// services/charts/chart_generator_factory.go
package services

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// ChartDataGeneratorFactory 用于根据图表类型和数据源类型返回对应的生成器
func ChartDataGeneratorFactory(chartType string, dataSourceType string, db *mongo.Database) (ChartDataGenerator, error) {
	// 如果是静态图表数据，直接返回静态数据生成器
	if dataSourceType == "static" {
		// 使用数据库连接创建 StaticChartDataGenerator 实例
		return NewStaticChartDataGenerator(db), nil
	}

	// 如果是动态图表数据，根据 chartType 返回不同的动态数据生成器
	switch chartType {
	case "score-trend":
		return &ScoreTrendChartDataGenerator{}, nil
	case "pass-rate":
		return &PassRateChartDataGenerator{}, nil
	default:
		return nil, fmt.Errorf("未找到图表类型: %s", chartType)
	}
}
