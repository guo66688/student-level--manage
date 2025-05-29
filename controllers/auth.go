// contollers/auth.go
package controllers

import (
	"fmt"
	"net/http"
	"time"

	"student-level-manage/config"
	"student-level-manage/middleware"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Login godoc
// @Summary 登录
// @Description 用户登录，获取 JWT token
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "登录请求"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	fmt.Println("🔥 Login 函数被调用")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "用户名或密码错误"})
		return
	}

	exp := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(middleware.JwtKey)

	// 设置 Cookie
	c.SetCookie("token", tokenStr, 3600*24, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"token": tokenStr, // 可用于 Postman 测试时直接拷贝
	})
}
