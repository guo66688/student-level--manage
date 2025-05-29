package middleware

import (
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// 权限校验中间件，传入 permission key
func RequirePermission(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user") // 在 JWT 中间件中设置的用户对象
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
			return
		}
		userID := user.(models.User).ID

		// 查询该用户所有角色
		var roleIDs []uint
		if err := config.DB.
			Table("user_roles").
			Select("role_id").
			Where("user_id = ?", userID).
			Scan(&roleIDs).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "查询角色失败"})
			return
		}

		// 查询这些角色是否包含所需权限
		var count int64
		err := config.DB.
			Table("role_permissions").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Where("role_permissions.role_id IN ? AND permissions.key = ?", roleIDs, key).
			Count(&count).Error

		if err != nil || count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "无权限访问"})
			return
		}

		c.Next()
	}
}
