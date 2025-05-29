package middleware

import (
	"net/http"
	"student-level-manage/config"

	"github.com/gin-gonic/gin"
)

// 权限验证中间件，参数为权限标识字符串（如："student:read"）
func RequirePermission(permissionKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户 ID（JWT 中间件应已写入）
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
			c.Abort()
			return
		}

		// 查询用户的权限列表
		var perms []string
		err := config.DB.
			Raw(`
				SELECT p.key FROM permissions p
				JOIN role_permissions rp ON rp.permission_id = p.id
				JOIN user_roles ur ON ur.role_id = rp.role_id
				WHERE ur.user_id = ?
			`, userID).
			Scan(&perms).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "权限查询失败"})
			c.Abort()
			return
		}

		// 判断是否包含目标权限
		has := false
		for _, p := range perms {
			if p == permissionKey {
				has = true
				break
			}
		}

		if !has {
			c.JSON(http.StatusForbidden, gin.H{"msg": "无权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
