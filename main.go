package main

import (
	"github.com/gin-gonic/gin"
	"CareerAnalysis/model"
	"encoding/json"
    "io/ioutil"
	"log"
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
	LoadConfig()
	r.Run(":" + config.Server.Port) // 监听并启动服务
}