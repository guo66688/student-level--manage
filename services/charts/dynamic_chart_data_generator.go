// services/charts/dynamic_chart_data_generator.go
package services

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
