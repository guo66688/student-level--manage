// routes.go
package routes

import (
	"fmt"
	"student-level-manage/controllers"
	"student-level-manage/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	fmt.Println("✅ 路由系统初始化成功")
	api := r.Group("/api")
	{
		api.POST("/auth/login", controllers.Login)

		// 加认证保护
		auth := api.Group("")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			auth.GET("/students", controllers.GetStudents)
			auth.POST("/students", controllers.AddStudent)

			auth.POST("/courses", controllers.AddCourse)
			auth.GET("/courses", controllers.GetCourses)

			auth.POST("/scores", controllers.AddScore)
			auth.GET("/scores", controllers.GetScores)

		}
	}

	return r
}
