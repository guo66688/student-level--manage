package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChartData struct {
	ID     primitive.ObjectID     `bson:"_id,omitempty"`
	Type   string                 `bson:"type"`
	Meta   map[string]interface{} `bson:"meta"`   // 可选元数据，例如 course_id 等
	Values interface{}            `bson:"values"` // 存储任意结构化数据，如统计图表数据
}
