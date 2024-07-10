package main

import (
	"github.com/gin-gonic/gin"
	"CareerAnalysis/model"
)

func main() {
	r := gin.Default()
	r.POST("/login", model.Login)
	r.POST("/signup", model.Signup)
	r.Run(":8886") // 监听并在 0.0.0.0:8886 上启动服务
}