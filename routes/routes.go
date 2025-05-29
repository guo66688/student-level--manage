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
			auth.PUT("/students/:id", controllers.UpdateStudent)
			auth.DELETE("/students/:id", controllers.DeleteStudent)
			auth.GET("/students/by_class", controllers.GetStudentsByClass)

			// 课程模块
			auth.POST("/courses", controllers.AddCourse)
			auth.GET("/courses", controllers.GetCourses)

			// 成绩模块
			auth.POST("/scores", controllers.AddScore)
			auth.GET("/scores", controllers.GetScores)

			// 班级模块
			auth.GET("/classes", controllers.GetClasses)
			auth.POST("/classes", controllers.AddClass)
			auth.PUT("/classes/:id", controllers.UpdateClass)
			auth.DELETE("/classes/:id", controllers.DeleteClass)

			// 数据分析模块
			auth.GET("/analytics/course-stats", controllers.GetCourseStats)
			auth.GET("/analysis/monthly", controllers.GetMonthlyStats)
			auth.GET("/analysis/pass_rate", controllers.GetPassRate)

			auth.GET("/analysis/rank", controllers.GetScoreRanking)
			auth.DELETE("/analysis/rank/cache", controllers.ClearScoreRankingCache)

			// 图表管理接口（Mongo 图表）
			chart := auth.Group("/charts")
			{
				chart.GET("", controllers.ListCharts)         // GET /api/charts?type=xxx
				chart.POST("", controllers.AddChart)          // 新增
				chart.PUT("/:id", controllers.UpdateChart)    // 编辑
				chart.DELETE("/:id", controllers.DeleteChart) // 删除（按 ID）

				chart.GET("/types", controllers.GetChartTypes)     // 图表类型列表
				chart.GET("/data", controllers.GetChartFromMongo)  // 前端按类型获取图表数据
				chart.DELETE("/data", controllers.DeleteChartData) // 按类型删除图表数据
			}
		}
	}

	return r
}
