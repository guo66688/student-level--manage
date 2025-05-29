// models/chart_data.go

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ChartData 图表数据（MongoDB）
type ChartData struct {
	ID     primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Type   string                 `bson:"type" json:"type" example:"monthly"`
	Meta   map[string]interface{} `bson:"meta" json:"meta"`
	Values interface{}            `bson:"values" json:"values"`
}
