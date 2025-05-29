package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChartData struct {
	ID     primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Type   string                 `bson:"type" json:"type"`
	Meta   map[string]interface{} `bson:"meta" json:"meta"`
	Values interface{}            `bson:"values" json:"values"`
}
