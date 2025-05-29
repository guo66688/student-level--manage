package models

type ChartData struct {
	Type   string                 `bson:"type"`   // 类型，比如 "monthly", "pass_rate"
	Meta   map[string]interface{} `bson:"meta"`   // 元数据，比如 {"course_id": 1}
	Values interface{}            `bson:"values"` // 实际图表数据（可以是任意结构）
}
