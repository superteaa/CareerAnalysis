package main

import "github.com/gin-gonic/gin"
import "CareerAnalysis/module"

func main() {
	r := gin.Default()
	r.GET("/login", module.Login)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}