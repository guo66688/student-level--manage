// services/charts/chart_data_generator.go
package services

// ChartDataGenerator 接口定义
// 每个图表类型都有一个生成器实现该接口
type ChartDataGenerator interface {
	GenerateData() ([]map[string]interface{}, error)              // 生成数据
	GenerateDataByID(id string) ([]map[string]interface{}, error) // 根据ID生成数据
}
