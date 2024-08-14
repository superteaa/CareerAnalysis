package model

import (
	"CareerAnalysis/baseClass"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Major 模型
type Major struct {
	ID   int // 专业id 1-通信工程，2-电子信息工程，3-计算机科学与技术，4-软件工程，5-网络工程，6-电子信息科学与技术，7-信息与计算科学
	Type int // 专业分类，0-工程类，1-信息类，2-理学类
	// Major_no  int  // 弃用
	Pic_Url string
	Intro   string
	// Com_User  string
	// Com_Title string
	// Com_Body  string
	// Star_Num  int
}

// Comment 模型
type Comment struct {
	ID        int // 评论id
	Major_id  int // 专业编号，1-通信工程，2-电子信息工程，3-计算机科学与技术，4-软件工程，5-网络工程，6-电子信息科学与技术，7-信息与计算科学
	Com_User  string
	Com_Title string
	Com_Body  string
	Star_Num  int
}

var MAJOR_TYPE_MAP = map[int]string{
	0: "工程类",
	1: "信息类",
	2: "理学类",
}

var MAJOR_NAME_MAP = map[int]string{
	1: "通信工程",
	2: "电子信息工程",
	3: "计算机科学与技术",
	4: "软件工程",
	5: "网络工程",
	6: "电子信息科学与技术",
	7: "信息与计算科学",
}

func GetMajorList(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	// 初始化数据库和Redis连接
	db := baseClass.GetDB()
	// rdb := baseClass.InitRedis()
	// defer rdb.Close()

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
			"major_id":      major.ID,
			"major_type":    MAJOR_TYPE_MAP[major.Type],
			"major_type_id": major.Type,
			"major_name":    MAJOR_NAME_MAP[major.ID],
		}
		result = append(result, majorListMap)
	}

	c.JSON(http.StatusOK, result)
}

func GetMajorDetail(c *gin.Context) {
	get_major_id := c.Query("major_id")
	if get_major_id == "" {
		c.JSON(http.StatusOK, gin.H{"error": "major_id不能为空"})
		return
	}
	major_id, err := strconv.Atoi(get_major_id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "参数格式错误"})
		log.Println(err.Error())
		return
	}

	db := baseClass.GetDB()
	var major Major
	var comments []Comment
	db_result := db.First(&major, major_id)

	if db_result.Error != nil {
		log.Println("查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	db_result = db.Where("major_id = ?", major_id).Find(&comments)

	if db_result.Error != nil {
		log.Println("查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	majorDetailMap := map[string]interface{}{
		"major_id": major.ID,
		"pic_url":  major.Pic_Url,
		"intro":    major.Intro,
		"name":     MAJOR_NAME_MAP[major.ID],
	}

	var commentListMap []map[string]interface{}
	for _, comment := range comments {
		commentMap := map[string]interface{}{
			"comment_id": comment.ID,
			"user":       comment.Com_User,
			"title":      comment.Com_Title,
			// "Body":  new.Body,
			"body": comment.Com_Body,
			"star": comment.Star_Num,
		}
		commentListMap = append(commentListMap, commentMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"major_info":  majorDetailMap,
		"commen_list": commentListMap,
	})
}
