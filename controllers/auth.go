// controllers/auth.go
package controllers

import (
	"net/http"
	"time"
	"log"
	"student-level-manage/config"
	"student-level-manage/middleware"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary 登录
// @Description 用户登录，获取 JWT token
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "登录请求" example({"username": "admin", "password": "123456"})
// @Success 200 {object} map[string]interface{} "返回 token 信息"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 401 {object} map[string]string "用户名或密码错误"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ JSON 解析失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误", "detail": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.
		Where("username = ?", req.Username).
		First(&user).Error; err != nil {
		log.Println("❌ 用户查询失败:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "用户名或密码错误", "detail": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		log.Println("❌ 密码校验失败:", err)
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
	tokenStr, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		log.Println("❌ Token 生成失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Token 签发失败", "detail": err.Error()})
		return
	}

	c.SetCookie("token", tokenStr, 3600*24, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"token": tokenStr,
	})
}