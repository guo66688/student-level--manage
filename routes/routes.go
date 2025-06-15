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
		// 登录接口
		api.POST("/auth/login", controllers.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			auth.GET("/doc/models", controllers.ShowAllModels)

			auth.GET("/users", controllers.ListUsers)
			auth.POST("/users", controllers.AddUser)
			auth.PUT("/users/:id", controllers.UpdateUser)
			auth.DELETE("/users/:id", controllers.DeleteUser)
			auth.PUT("/users/:id/reset_password", controllers.ResetPassword)
			auth.GET("/user/roles", controllers.GetUserRoles) // 用于前端获取当前登录用户角色（新文件 user.go 中的）
			auth.GET("/user/permissions", controllers.GetUserPermissions)

			// ✅ 权限管理
			auth.POST("/permissions", controllers.AddPermission)
			auth.GET("/permissions", controllers.GetPermissions)
			auth.DELETE("/permissions/:id", controllers.DeletePermission)

			// ✅ 角色管理
			auth.POST("/roles", controllers.AddRole)
			auth.GET("/roles", controllers.GetRoles)
			auth.DELETE("/roles/:id", controllers.DeleteRole)
			auth.POST("/role/permissions", controllers.SetRolePermissions)
			auth.GET("/role/permission_ids", controllers.GetPermissionsByRole)

			// ✅ 用户角色绑定
			auth.POST("/user/role", controllers.SetUserRole)
			auth.GET("/user/role", controllers.GetUserRole) // 用于后台指定用户查角色
			// ✅ 学生管理
			auth.GET("/students", middleware.RequirePermission("student:view"), controllers.GetStudents)
			auth.POST("/students", middleware.RequirePermission("student:add"), controllers.AddStudent)
			auth.PUT("/students/:id", middleware.RequirePermission("student:update"), controllers.UpdateStudent)
			auth.DELETE("/students/:id", middleware.RequirePermission("student:delete"), controllers.DeleteStudent)
			auth.GET("/students/by_class", middleware.RequirePermission("student:view"), controllers.GetStudentsByClass)

			// ✅ 课程模块
			auth.GET("/courses", middleware.RequirePermission("course:view"), controllers.GetCourses)
			auth.POST("/courses", middleware.RequirePermission("course:add"), controllers.AddCourse)
			auth.PUT("/courses/:id", middleware.RequirePermission("course:update"), controllers.UpdateCourse)
			auth.DELETE("/courses/:id", middleware.RequirePermission("course:delete"), controllers.DeleteCourse)

			// ✅ 成绩管理

			auth.POST("/scores", middleware.RequirePermission("score:add"), controllers.AddScore)
			auth.GET("/scores", middleware.RequirePermission("score:view"), controllers.GetScores)
			auth.PUT("/scores/:id", middleware.RequirePermission("score:update"), controllers.UpdateScore)
			auth.DELETE("/scores/:id", middleware.RequirePermission("score:delete"), controllers.DeleteScore)

			// ✅ 班级管理
			auth.GET("/classes", middleware.RequirePermission("class:view"), controllers.GetClasses)
			auth.POST("/classes", middleware.RequirePermission("class:add"), controllers.AddClass)
			auth.PUT("/classes/:id", middleware.RequirePermission("class:update"), controllers.UpdateClass)
			auth.DELETE("/classes/:id", middleware.RequirePermission("class:delete"), controllers.DeleteClass)

			// ✅ 数据分析（无需权限控制）
			auth.GET("/analytics/course-stats", controllers.GetCourseStats)
			auth.GET("/analysis/class_avg", controllers.GetClassAvgScore)
			auth.GET("/analysis/monthly", controllers.GetMonthlyStats)
			auth.GET("/analysis/pass_rate", controllers.GetPassRate)
			auth.GET("/analysis/exam_count", controllers.GetExamCount)
			auth.GET("/dashboard/stats", controllers.GetDashboardStats)
			auth.GET("/analysis/rank", controllers.GetScoreRanking)
			auth.DELETE("/analysis/rank/cache", controllers.ClearScoreRankingCache)

			// ✅ 图表管理（MongoDB）
			charts := auth.Group("/charts")
			{
				charts.GET("", controllers.ListCharts)
				charts.POST("", middleware.RequirePermission("chart:add"), controllers.AddChart)
				charts.PUT("/:id", middleware.RequirePermission("chart:update"), controllers.UpdateChart)
				charts.DELETE("/:id", middleware.RequirePermission("chart:delete"), controllers.DeleteChart)
				charts.GET("/types", controllers.GetChartTypes)
				charts.GET("/data", controllers.GetChartFromMongo)
				charts.DELETE("/data", controllers.DeleteChartData)
			}
		}
	}
	

	return r
}
