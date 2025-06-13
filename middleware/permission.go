// middleware/permission.go
package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// 获取用户权限列表，先尝试从 Redis 缓存中获取，若没有则从数据库中查询并缓存
func getUserPermissions(userID uint) ([]string, error) {
	var names []string
	cacheKey := fmt.Sprintf("user:%d:permissions", userID)
	val, err := config.Redis.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，返回权限列表
		json.Unmarshal([]byte(val), &names)
		return names, nil
	}

	// 缓存未命中，从数据库查询权限
	var roleIDs []uint
	config.DB.Table("user_roles").Select("role_id").Where("user_id = ?", userID).Scan(&roleIDs)
	config.DB.
		Table("role_permissions").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Pluck("permissions.name", &names)

	// 缓存查询结果
	data, _ := json.Marshal(names)
	config.Redis.Set(config.Ctx, cacheKey, data, time.Hour) // 缓存 1 小时
	return names, nil
}

// 权限校验中间件，传入 permission name
func RequirePermission(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
			return
		}
		userID := user.(models.User).ID

		// 获取用户权限
		names, err := getUserPermissions(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "获取权限失败"})
			return
		}

		// 打印当前用户权限列表，便于调试
		fmt.Println("当前用户权限列表:", names)

		// 校验是否有该权限
		permissionMap := make(map[string]struct{})
		for _, n := range names {
			permissionMap[n] = struct{}{}
		}

		if _, found := permissionMap[name]; !found {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "无权限访问"})
			return
		}

		// 权限验证通过，继续处理请求
		c.Next()
	}
}
