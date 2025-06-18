// services/charts/dynamic_chart_data_generator.go
package services

import "fmt"

// ScoreTrendChartDataGenerator 用于生成成绩趋势图表的数据
type ScoreTrendChartDataGenerator struct{}

// 实现 ChartDataGenerator 接口
func (g *ScoreTrendChartDataGenerator) GenerateData() ([]map[string]interface{}, error) {
	var values []map[string]interface{}

	// 假设从数据库中查询成绩数据并生成数据
	// 这里我们手动插入一些数据示例
	values = append(values, map[string]interface{}{"x": "2024-01", "y": 85.3})
	values = append(values, map[string]interface{}{"x": "2024-02", "y": 88.5})
	values = append(values, map[string]interface{}{"x": "2024-03", "y": 90.0})

	return values, nil
}

// 根据ID获取数据的方法
func (g *ScoreTrendChartDataGenerator) GenerateDataByID(id string) ([]map[string]interface{}, error) {
	// 假设根据ID返回不同的数据
	if id == "6852592bdc19f14fc8d8957c" {
		// 返回特定图表数据
		return []map[string]interface{}{
			{"x": "2024-01", "y": 75.3},
			{"x": "2024-02", "y": 82.0},
			{"x": "2024-03", "y": 88.5},
		}, nil
	}
	// 如果没有找到对应的ID，返回错误
	return nil, fmt.Errorf("未找到ID为 %s 的数据", id)
}

// PassRateChartDataGenerator 用于生成通过率图表的数据
type PassRateChartDataGenerator struct{}

// 实现 ChartDataGenerator 接口
func (g *PassRateChartDataGenerator) GenerateData() ([]map[string]interface{}, error) {
	var values []map[string]interface{}

	// 假设从数据库中查询每门课的通过率数据
	// 这里我们手动插入一些数据示例
	values = append(values, map[string]interface{}{"x": "Course 1", "y": 92.5})
	values = append(values, map[string]interface{}{"x": "Course 2", "y": 85.0})

	return values, nil
}

// 根据ID获取数据的方法
func (g *PassRateChartDataGenerator) GenerateDataByID(id string) ([]map[string]interface{}, error) {
	// 根据ID返回特定的数据
	if id == "123456789" {
		return []map[string]interface{}{
			{"x": "Course 1", "y": 91.0},
			{"x": "Course 2", "y": 84.5},
		}, nil
	}
	// 如果没有找到对应的ID，返回错误
	return nil, fmt.Errorf("未找到ID为 %s 的数据", id)
}
