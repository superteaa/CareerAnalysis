package model

import (
	"CareerAnalysis/baseClass"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// User 模型
type User struct {
	ID       int `gorm:"primaryKey"`
	Username string
	Password string
	Email    string `gorm:"uniqueIndex"`
	Token    string
	// Avatar   string
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
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid request"})
		return
	}

	// 初始化数据库和Redis连接
	db := baseClass.GetDB()

	var user User
	if err := db.Where("username = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Username does not exist"})
		return
	}

	if !CheckPassword(request.Password, user.Password) {
		c.JSON(http.StatusOK, gin.H{"error": "Username and password do not match"})
		return
	}

	avatarURL := fmt.Sprintf("/uploads/%s.jpg", user.Email)

	if user.Token == "" {
		// 生成JWT会话令牌
		sessionToken, err := baseClass.GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_token": sessionToken, "username": user.Username, "avatar": avatarURL})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_token": user.Token, "username": user.Username, "avatar": avatarURL})

}

// Signup 处理用户注册请求
func Signup(c *gin.Context) {
	// 解析表单数据
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	// captchaId := c.PostForm("captchaId")
	// value := c.PostForm("value")

	// // 验证验证码
	// if !Verifycaptcha(captchaId, value) {
	// 	c.JSON(http.StatusOK, gin.H{"error": "Invalid captcha"})
	// 	return
	// }

	// 处理头像上传
	file, err := c.FormFile("avatar")
	if err != nil {
		log.Println("signup: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar upload failed"})
		return
	}

	// 保存头像到服务器（假设保存到 "uploads/" 目录）
	avatarPath := fmt.Sprintf("/uploads/%s.jpg", email)
	if err := c.SaveUploadedFile(file, "."+avatarPath); err != nil {
		log.Println("signup: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save avatar"})
		return
	}

	// 初始化数据库连接
	db := baseClass.GetDB()

	// 检查用户名是否已经存在
	var existingUser User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"error": "Email already exists"})
		return
	}

	// 创建新用户
	newUser := User{
		Username: username,
		Password: password,
		Email:    email,
		// Avatar:   avatarPath,
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
	c.JSON(http.StatusOK, gin.H{"message": "User signup successfully", "session_token": sessionToken, "username": newUser.Username, "avatar": avatarPath})
}
