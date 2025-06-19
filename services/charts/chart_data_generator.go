// services/chart_data_generator.go
package services

// ChartDataGenerator 用于生成图表数据的接口
type ChartDataGenerator interface {
    // GenerateData 通过某些条件生成数据
    GenerateData() ([]map[string]interface{}, error)
    // 根据ID获取数据的方法
    GenerateDataByID(id string) ([]map[string]interface{}, error)
}
