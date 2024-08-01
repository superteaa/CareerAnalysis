package model

import (
	"CareerAnalysis/baseClass"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

// new 模型
type New struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	Body     string
	Date     string
	Icon_url string
}

func GetNewList(c *gin.Context) {
	// 初始化数据库和Redis连接
	db := baseClass.InitDB()
	rdb := baseClass.InitRedis()
	defer rdb.Close()

	var news []New
	db_result := db.Find(&news)
	if db_result.Error != nil {
		log.Println("查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	var result []map[string]interface{}
	for _, new := range news {
		newMap := map[string]interface{}{
			"news_id": new.ID,
			"title":   new.Title,
			// "Body":  new.Body,
			"date":     new.Date,
			"icon_url": new.Icon_url,
		}
		result = append(result, newMap)
	}

	c.JSON(http.StatusOK, result)
}

func GetNews(c *gin.Context) {
	news_id := c.Query("news_id")
	// 初始化数据库和Redis连接
	db := baseClass.InitDB()
	var news New
	db_result := db.Where("id = ?", news_id).First(&news)
	if db_result.Error != nil {
		log.Println("查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	result := map[string]interface{}{
		"news_id":  news.ID,
		"title":    news.Title,
		"body":     news.Body,
		"date":     news.Date,
		"icon_url": news.Icon_url,
	}

	c.JSON(http.StatusOK, result)

}
