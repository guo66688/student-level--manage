package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"
	"time"

	"github.com/gin-gonic/gin"
)

// 权限校验中间件，传入 permission key

func RequirePermission(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
			return
		}
		userID := user.(models.User).ID

		var keys []string
		cacheKey := fmt.Sprintf("user:%d:permissions", userID)
		val, err := config.Redis.Get(config.Ctx, cacheKey).Result()
		if err == nil {
			json.Unmarshal([]byte(val), &keys)
		} else {
			// 查询数据库
			var roleIDs []uint
			config.DB.Table("user_roles").Select("role_id").Where("user_id = ?", userID).Scan(&roleIDs)
			config.DB.
				Table("role_permissions").
				Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
				Where("role_permissions.role_id IN ?", roleIDs).
				Pluck("permissions.key", &keys)

			// 缓存权限
			data, _ := json.Marshal(keys)
			config.Redis.Set(config.Ctx, cacheKey, data, time.Hour)
		}

		// 校验是否有该权限
		found := false
		for _, k := range keys {
			if k == key {
				found = true
				break
			}
		}
		if !found {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "无权限访问"})
			return
		}
		c.Next()
	}
}
