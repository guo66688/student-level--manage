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
			// 学生模块
			auth.GET("/students", controllers.GetStudents)
			auth.POST("/students", controllers.AddStudent)

			// 课程模块
			auth.POST("/courses", controllers.AddCourse)
			auth.GET("/courses", controllers.GetCourses)

			// 成绩模块
			auth.POST("/scores", controllers.AddScore)
			auth.GET("/scores", controllers.GetScores)

			// 数据分析模块
			auth.GET("/analytics/course-stats", controllers.GetCourseStats)

			auth.GET("/analysis/rank", controllers.GetScoreRanking)
			auth.GET("/analysis/monthly", controllers.GetMonthlyStats)
			auth.GET("/analysis/pass_rate", controllers.GetPassRate)

			auth.GET("/analysis/rank_redis", controllers.GetScoreRankingRedis)

		}
	}

	return r
}
