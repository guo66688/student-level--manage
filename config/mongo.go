package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("❌ MongoDB连接失败:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB Ping失败:", err)
	}

	MongoDB = client
	log.Println("✅ MongoDB连接成功")
}
