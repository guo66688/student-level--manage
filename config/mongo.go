package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoCtx context.Context

func InitMongo() {
	MongoCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://admin:123456@localhost:27017")

	var err error
	MongoClient, err = mongo.Connect(MongoCtx, clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB连接失败:", err)
	}
}
