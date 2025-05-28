package routes

import (
	"student-level-manage/controllers"
	"student-level-manage/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/auth/login", controllers.Login)

		// 加认证保护
		auth := api.Group("")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			auth.GET("/students", controllers.GetStudents)
			auth.POST("/students", controllers.AddStudent)
		}
	}

	return r
}
