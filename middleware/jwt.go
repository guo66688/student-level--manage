// middleware/jwt.go
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("secret_key")

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// 允许从 Cookie 中读取 token
			cookieToken, err := c.Cookie("token")
			if err == nil {
				tokenStr = cookieToken
			}
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "未提供有效 token"})
			return
		}

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

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
