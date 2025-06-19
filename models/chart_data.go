// models/chart_data.go

package models

// Meta 结构体，表示图表元数据
type Meta struct {
	Key   string `bson:"Key" json:"Key"`
	Value string `bson:"Value" json:"Value"`
}

// ChartValues 表示一个图表的 x 和 y 数据
type ChartValues struct {
	X []string  `bson:"x" json:"x"` // x 数据，字符串数组
	Y []float64 `bson:"y" json:"y"` // y 数据，浮动数字数组
}

// ChartData 图表数据（MongoDB）
type ChartData struct {
	ID             string        `bson:"_id" json:"id"`                            // 图表 ID
	Type           string        `bson:"type" json:"type"`                         // 图表类型
	Meta           []Meta        `bson:"meta" json:"meta"`                         // 元数据，使用数组形式
	Values         []ChartValues `bson:"values" json:"values"`                     // 图表数据，包含 X 和 Y（数组类型）
	DataSourceType string        `bson:"data_source_type" json:"data_source_type"` // 数据源类型
}
