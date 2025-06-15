// @title 学生成绩管理系统 API
// @version 1.1
// @description 用于管理学生、课程、成绩和图表数据的后台 API。
// @contact.name icoffee
// @contact.email 1596740959@qq.com
// @host localhost:8080
// @BasePath /api

// main.go
package main

import (
	"log"
	"time"

	"student-level-manage/config"
	"student-level-manage/models"
	"student-level-manage/redisop"
	"student-level-manage/routes"

	swaggerFiles "github.com/swaggo/files"

	// ✅ 替代旧的 swaggerFiles
	_ "student-level-manage/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
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
		&models.Class{},
	)

	// 启动路由服务
	r := routes.InitRouter()

	// ✅ 注册 Swagger 接口文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run("0.0.0.0:8080")
	log.Println("🔥 Air Reload Test")

}
