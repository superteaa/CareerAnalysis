package main

import (
	"CareerAnalysis/baseClass"
	"CareerAnalysis/model"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

// LoadConfig 从文件中加载配置
func LoadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func main() {

	// 打开日志文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法打开日志文件:", err)
	}
	// 将日志输出设置到文件
	log.SetOutput(file)

	log.Println("应用程序启动")

	baseClass.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://8.138.115.27:8866"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	r.POST("/login", model.Login)
	// r.POST("/signup", model.Signup)
	r.POST("/upload_avatar")
	// r.GET("/captcha/:captchaId", model.Getcaptchaimg)
	// r.GET("/captcha", model.Createcaptchaid)
	authorized := r.Group("/")
	authorized.Use(baseClass.ValidateJWT())
	{
		// question_r := authorized.Group("/question")
		// {
		// 	// question_r.GET("/roles", model.SelectRole)
		// 	question_r.POST("/submit", model.DealQ)
		// }
		job_r := authorized.Group("/job")
		{
			job_r.GET("/get-list", model.GetJobList)
		}
		study_r := authorized.Group("/study")
		{
			study_r.POST("/add-plan", model.AddPlan)
			study_r.GET("/get-data", model.GetStudyData)
			study_r.GET("/get-plan-list", model.GetPlanList)
			study_r.GET("/get-plan-detail", model.GetPlanDetail)
			study_r.POST("/change-plan", model.ChangePlan)
			study_r.GET("/get-subject-map", model.GetSubjectMap)
			study_r.GET("/get-learned-map", model.GetSkillTree)
		}
		new_r := authorized.Group("/news")
		{
			new_r.GET("/get-list", model.GetNewList)
			new_r.GET("/get-detail", model.GetNews)
		}
		dataAnalys_r := authorized.Group("/data")
		{
			dataAnalys_r.GET("/get-detail", model.GetSubjectRate)
		}
		major_r := authorized.Group("/major")
		{
			major_r.GET("/get-list", model.GetMajorList)
			major_r.GET("/get-detail", model.GetMajorDetail)
		}
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current directory:", dir)
	r.Run(":" + config.Server.Port) // 监听并启动服务
}
