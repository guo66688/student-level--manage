package main

import (
	"log"
	"student-level-manage/config"
	"student-level-manage/models"
	"student-level-manage/redisop"
	"student-level-manage/routes"
	"time"
)

func main() {
	// 初始化数据库连接
	config.InitDB()
	config.InitRedis()
	config.InitMongo()

	// 定时刷新 Redis 排行榜
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			if err := redisop.RefreshAllRanksFromDB(); err != nil {
				log.Println("❌ 刷新 Redis 排行榜失败:", err)
			} else {
				log.Println("✅ Redis 排行榜已刷新")
			}
			<-ticker.C
		}
	}()

	// 自动创建表
	config.DB.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Course{},
		&models.Score{},
		// ✅ 新增权限管理相关表
		&models.Permission{},
		&models.Role{},
		&models.RolePermission{},
		&models.UserRole{},
	)

	// 启动路由服务
	r := routes.InitRouter()
	r.Run(":8080")
}
