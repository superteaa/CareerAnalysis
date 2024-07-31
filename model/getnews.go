package model

import (
	"CareerAnalysis/baseClass"

	"net/http"

	"github.com/gin-gonic/gin"
)

// new 模型
type New struct {
	ID    uint `gorm:"primaryKey"`
	Title string
	Body  string
	Date  string
	Icon  string
}

func GetNewList(c *gin.Context) {
	// 初始化数据库和Redis连接
	db := baseClass.InitDB()
	rdb := baseClass.InitRedis()
	defer rdb.Close()

	var news []New
	db.Find(&news)

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
