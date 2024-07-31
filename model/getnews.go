package model

import (
	"CareerAnalysis/baseClass"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

// new 模型
type New struct {
	ID    uint `gorm:"primaryKey"`
	Title string
	Body  string
	Date  string
	Icon  string `gorm:"column:icon_url"`
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
			"Title": new.Title,
			"Body":  new.Body,
			"Date":  new.Date,
			"Icon":  new.Icon,
		}
		result = append(result, newMap)
	}

	c.JSON(http.StatusOK, result)
}
