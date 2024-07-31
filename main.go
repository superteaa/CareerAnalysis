package main

import (
	"CareerAnalysis/baseClass"
	"CareerAnalysis/model"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

var config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

// LoadConfig 从文件中加载配置
func LoadConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func main() {
	r := gin.Default()
	r.POST("/login", model.Login)
	r.POST("/signup", model.Signup)
	r.GET("/captcha/:captchaId", model.Getcaptchaimg)
	r.GET("/captcha", model.Createcaptchaid)
	authorized := r.Group("/")
	authorized.Use(baseClass.ValidateJWT())
	{
		authorized.GET("/getnewlist", model.GetNewList)
	}
	// r.POST("/hello", func(c *gin.Context) {
	// 	if err := baseClass.ValidateJWT(c); err == nil {
	// 		ID, exists := c.Get("userID")
	// 		if exists {
	// 			c.JSON(200, gin.H{
	// 				"ID": ID,
	// 			})
	// 		} else {
	// 			c.JSON(200, gin.H{
	// 				"ID": ID,
	// 			})
	// 		}
	// 	}

	// })
	LoadConfig()
	r.Run(":" + config.Server.Port) // 监听并启动服务
}
