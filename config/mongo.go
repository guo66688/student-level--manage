// config\mongo.go
package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client   // 客户端连接
	MongoDB     *mongo.Database // 数据库对象
)

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 更新为使用认证的MongoDB连接
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:123456@localhost:27017").
		SetAuth(options.Credential{
			Username:   "admin",
			Password:   "123456",
			AuthSource: "admin", // 认证数据库，通常是 admin
		}))
	if err != nil {
		log.Fatal("❌ MongoDB连接失败:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB Ping失败:", err)
	}

	MongoClient = client
	MongoDB = MongoClient.Database("analysis") // ✅ 指定数据库名
	log.Println("✅ MongoDB连接成功")
}
