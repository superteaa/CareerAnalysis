package model

import (
	"CareerAnalysis/baseClass"
	"fmt"
	"log"
	"net/http"
	"os"

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
	Is_new int
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
	if err := db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "User does not exist"})
		return
	}

	if !CheckPassword(request.Password, user.Password) {
		c.JSON(http.StatusOK, gin.H{"error": "Email and password do not match"})
		return
	}

	// 支持的图片类型
	supportedExtensions := []string{".jpg", ".png", ".gif"}
	var avatarURL string

	// 遍历每种后缀，查找是否存在相应的文件
	for _, ext := range supportedExtensions {
		filePath := fmt.Sprintf("./uploads/%s%s", request.Email, ext)
		if _, err := os.Stat(filePath); err == nil {
			avatarURL = fmt.Sprintf("/uploads/%s%s", request.Email, ext)
			break
		}
	}

	// 如果没有找到头像文件，可以返回一个默认头像 URL
	if avatarURL == "" {
		avatarURL = "/uploads/default-avatar.png" // 默认头像路径
	}

	var is_test int

	record := map[string]interface{}{}
	db_result := db.Table("recomment_jobs").Where("user_id = ?", user.ID).Find(&record)
	if db_result.RowsAffected == 0 {
		is_test = 2
	} else {
		is_test = 1
	}

	if user.Token == "" {
		// 生成JWT会话令牌
		sessionToken, err := baseClass.GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_token": sessionToken, "username": user.Username, "avatar": avatarURL, "is_test": is_test, "is_new": user.Is_new})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_token": user.Token, "username": user.Username, "avatar": avatarURL, "is_test": is_test, "is_new": user.Is_new})

}

// Signup 处理用户注册请求
func Signup(c *gin.Context) {
	// 解析表单数据
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	captchaId := c.PostForm("captchaId")
	value := c.PostForm("value")

	// 验证验证码
	if !Verifycaptcha(captchaId, value) {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid captcha"})
		return
	}

	var avatarURL string
	files := c.Request.MultipartForm.File["files"]
	if len(files) != 0 {
		file := files[0]
		// 打开文件并读取文件头部信息
		src, err := file.Open()
		if err != nil {
			log.Println("Failed to open uploaded file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
			return
		}
		defer src.Close()

		// 读取文件的 MIME 类型
		buffer := make([]byte, 512)
		if _, err := src.Read(buffer); err != nil {
			log.Println("Failed to read file header:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
			return
		}
		fileType := http.DetectContentType(buffer)

		// 根据 MIME 类型设置文件后缀
		var fileExt string
		switch fileType {
		case "image/jpeg":
			fileExt = ".jpg"
		case "image/png":
			fileExt = ".png"
		case "image/gif":
			fileExt = ".gif"
		default:
			log.Println("Unsupported file type:", fileType)
			fileExt = ".png"
		}

		// 回到文件开头读取全部数据
		if _, err := src.Seek(0, 0); err != nil {
			log.Println("Failed to seek file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
			return
		}

		// 保存头像到服务器（假设保存到 "uploads/" 目录）
		avatarPath := fmt.Sprintf("./uploads/%s%s", email, fileExt)
		avatarURL = fmt.Sprintf("/uploads/%s%s", email, fileExt)
		if err := c.SaveUploadedFile(file, avatarPath); err != nil {
			log.Println("Failed to save avatar:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save avatar"})
			return
		}
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
		Is_new: 1,
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
	c.JSON(http.StatusOK, gin.H{"message": "User signup successfully", "session_token": sessionToken, "username": newUser.Username, "avatar": avatarURL})
}

func UpdateIsNew(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("UpdateIsNew:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	var user User
	user.ID = int(userID.(uint32))
	// 初始化数据库连接
	db := baseClass.GetDB()
	db.Model(&user).Update("is_new", 0)
	log.Println("新用户完成引导userID:", user.ID)
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}
