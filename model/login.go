package model

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "CareerAnalysis/baseClass"
    "fmt"
)

// User 模型
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex"`
    Password string
    Email string
}

func (User) TableName() string {
    return "user"
}

func CheckPassword(r_password string, u_password string) bool {
    if r_password == u_password {
        return true;
    } else {
        return false;
    }
}

// login 处理登录请求
func Login(c *gin.Context) {
    var request struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // 初始化数据库和Redis连接
    db := baseClass.InitDB()
    rdb := baseClass.InitRedis()
    defer rdb.Close()

    var user User
    if err := db.Where("username = ?", request.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    if !CheckPassword(request.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // 生成JWT会话令牌
    sessionToken, err := baseClass.GenerateJWT(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_token": sessionToken})
}

// Signup 处理用户注册请求
func Signup(c *gin.Context) {
    var request struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
        Email string `json:"email" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // 初始化数据库连接
    db := baseClass.InitDB()

    // 检查用户名是否已经存在
    var existingUser User
    if err := db.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
        return
    }

    // 创建新用户
    newUser := User{
        Username: request.Username,
        Password: request.Password,
        Email: request.Email,
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

    // 返回成功响应
    c.JSON(http.StatusOK, gin.H{"message": "User signup successfully", "session_token": sessionToken})
}
