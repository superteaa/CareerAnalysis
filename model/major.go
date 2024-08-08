package model

import (
	"CareerAnalysis/baseClass"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Major 模型
type Major struct {
	ID       uint // 专业id
	Type     int  // 专业分类，0-工程类，1-信息类，2-理学类
	Major_no int  // 专业编号，0-通信工程，1-电子信息工程，2-计算机科学与技术，3-软件工程，4-网络工程，5-电子信息科学与技术，6-信息与计算科学
}

var MAJOR_TYPE_MAP = map[int]string{
	0: "工程类",
	1: "信息类",
	2: "理学类",
}

var MAJOR_NAME_MAP = map[int]string{
	0: "通信工程",
	1: "电子信息工程",
	2: "计算机科学与技术",
	3: "软件工程",
	4: "网络工程",
	5: "电子信息科学与技术",
	6: "信息与计算科学",
}

func GetMajorList(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	// 初始化数据库和Redis连接
	db := baseClass.InitDB()
	rdb := baseClass.InitRedis()
	defer rdb.Close()

	var majors []Major
	db_result := db.Find(&majors)
	if db_result.Error != nil {
		log.Println("查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	var result []map[string]interface{}
	for _, major := range majors {
		majorListMap := map[string]interface{}{
			"major_id":   major.ID,
			"major_type": MAJOR_TYPE_MAP[major.Type],
			"major_name": MAJOR_TYPE_MAP[major.Major_no],
		}
		result = append(result, majorListMap)
	}

	c.JSON(http.StatusOK, result)
}
