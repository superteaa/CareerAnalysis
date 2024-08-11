package model

import (
	"CareerAnalysis/baseClass"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// User 模型
type User struct {
	ID       int    `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Password string
	Email    string
	Token    string
}

func (User) TableName() string {
	return "user"
}

func CheckPassword(r_password string, u_password string) bool {
	if r_password == u_password {
		return true
	} else {
		return false
	}
}

// login 处理登录请求
func Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid request"})
		return
	}

	// 初始化数据库和Redis连接
	db := baseClass.InitDB()

	var user User
	if err := db.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Username does not exist"})
		return
	}

	if !CheckPassword(request.Password, user.Password) {
		c.JSON(http.StatusOK, gin.H{"error": "Username and password do not match"})
		return
	}

	if user.Token == "" {
		// 生成JWT会话令牌
		sessionToken, err := baseClass.GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_token": sessionToken})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_token": user.Token})

}

// Signup 处理用户注册请求
func Signup(c *gin.Context) {
	var request struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Email     string `json:"email" binding:"required"`
		CaptchaId string `json:"captchaId" binding:"required"`
		Value     string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "Invalid request"})
		return
	}

	if !Verifycaptcha(request.CaptchaId, request.Value) {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid captcha"})
		return
	}

	// 初始化数据库连接
	db := baseClass.InitDB()

	// 检查用户名是否已经存在
	var existingUser User
	if err := db.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"error": "Username already exists"})
		return
	}

	// 创建新用户
	newUser := User{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 生成JWT会话令牌
	sessionToken, err := baseClass.GenerateJWT(newUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	db.Model(&newUser).Update("token", sessionToken)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "User signup successfully", "session_token": sessionToken})
}
