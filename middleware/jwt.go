// middleware/jwt.go
package middleware

import (
	"net/http"
	"strings"
	"time"

	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("secret_key")

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

// JWT 鉴权中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 token 字符串
		var tokenStr string
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else if cookie, err := c.Cookie("token"); err == nil {
			tokenStr = cookie
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "未提供有效 token"})
			return
		}

		// 2. 解析并校验 token
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "token 无效或过期"})
			return
		}

		claims := token.Claims.(*Claims)
		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "token 已过期"})
			return
		}

		// 3. 根据 userID 从数据库加载用户记录
		var user models.User
		if err := config.DB.First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "用户不存在"})
			return
		}

		// 4. 将用户信息存入 context，供后续中间件/handler 使用
		c.Set("user", user)
		c.Set("user_id", user.ID)

		c.Next()
	}
}
