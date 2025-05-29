package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client // ✅ 修改名为 Redis
var Ctx = context.Background()

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 默认无密码
		DB:       0,
	})
	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("❌ Redis连接失败:", err)
	}
}
