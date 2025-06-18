// services/static_chart_data_generator.go
package services

// StaticChartDataGenerator 用于生成静态图表的数据
type StaticChartDataGenerator struct{}

// 实现 ChartDataGenerator 接口
func (g *StaticChartDataGenerator) GenerateData() ([]map[string]interface{}, error) {
	var values []map[string]interface{}

	// 假设静态数据直接返回
	values = append(values, map[string]interface{}{"x": "2024-01", "y": 85.3})
	values = append(values, map[string]interface{}{"x": "2024-02", "y": 88.5})
	values = append(values, map[string]interface{}{"x": "2024-03", "y": 90.0})

	return values, nil
}
