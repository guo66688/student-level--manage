// config\mongo.go
package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client // 客户端连接
var MongoDB *mongo.Database   // 数据库对象

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

	MongoClient = client
	MongoDB = MongoClient.Database("student_manage") // ✅ 指定数据库名
	log.Println("✅ MongoDB连接成功")
}
