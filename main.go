// @title 学生成绩管理系统 API
// @version 1.0
// @description 用于管理学生、课程、成绩和图表数据的后台 API。
// @contact.name 开发者
// @contact.email dev@example.com
// @host localhost:8080
// @BasePath /api

// main.go
package main

import (
	"log"
	"student-level-manage/config"
	"student-level-manage/models"
	"student-level-manage/redisop"
	"student-level-manage/routes"
	"time"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "student-level-manage/docs"
)

// swag 会生成这个包

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
		&models.Permission{},
		&models.Role{},
		&models.RolePermission{},
		&models.UserRole{},
	)

	// 启动路由服务
	r := routes.InitRouter()

	// ✅ 注册 Swagger 接口文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
