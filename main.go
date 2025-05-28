package main

import (
	"student-level-manage/config"
	"student-level-manage/models"
	"student-level-manage/routes"
)

func main() {
	// 初始化数据库连接
	config.InitDB()
	config.InitRedis()
	config.InitMongo()

	// 自动创建表
	config.DB.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Course{},
		&models.Score{},
	)

	// 启动路由服务
	r := routes.InitRouter()
	r.Run(":8080")
}
