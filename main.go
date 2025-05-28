package main

import (
	"student-system/config"
	"student-system/models"
	"student-system/routes"
)

func main() {
	config.InitDB()
	config.InitRedis()
	config.InitMongo()

	config.DB.AutoMigrate(&models.User{}, &models.Student{}, &models.Course{}, &models.Score{})

	r := routes.InitRouter()
	r.Run(":8080")
}
