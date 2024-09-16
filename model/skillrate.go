package model

import (
	"CareerAnalysis/baseClass"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Job 模型
type Job struct {
	ID           int
	Type         int // 工作分类，0-工程类，1-信息类，2-理学类
	Main_skill   string
	Expand_skill string `gorm:"column:expand_skill"`
	Data_rows    int
}

// SubjectRate 模型
type SubjectRate struct {
	ID         int `gorm:"primaryKey"`
	Subject_Id int
	Skill_Name string
	Major_id   int
	Rate       float32
	Study_url  string
}

var JOB_TYPE_MAP = map[int]string{
	0: "工程类",
	1: "信息类",
	2: "理学类",
}

var JOB_NAME_MAP = map[int]string{
	1: "产品经理",
	2: "测试工程师",
	3: "后端工程师",
	4: "前端工程师",
	5: "算法工程师",
}

func (SubjectRate) TableName() string {
	return "skills"
}

func GetSubjectRate(c *gin.Context) {
	// major_id := c.Query("major_id")
	job_id := c.Query("job_id")

	db := baseClass.GetDB()

	var subject_rate []SubjectRate
	db_result := db.Where("job_id = ?", job_id).Find(&subject_rate)
	if db_result.Error != nil {
		log.Println("SkillRate查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	// record_rows := db_result.RowsAffected

	var subjects_info []map[string]interface{}

	for _, v := range subject_rate {

		sigle := map[string]interface{}{
			"value":        v.Rate,
			"subject_name": v.Skill_Name,
		}
		if v.Study_url != "" {
			log.Println(v.Study_url)
			sigle["study_url"] = v.Study_url
		}
		subjects_info = append(subjects_info, sigle)
	}

	var job_dec Job
	db_result = db.Where("id = ?", job_id).First(&job_dec)
	if db_result.Error != nil {
		log.Println("Job查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	now := time.Now()
	threeDaysAgo := now.AddDate(0, 0, -3)
	threeDaysAgoMidnight := time.Date(threeDaysAgo.Year(), threeDaysAgo.Month(), threeDaysAgo.Day(), 0, 0, 0, 0, threeDaysAgo.Location()).Unix()
	result := map[string]interface{}{
		"subject_value": subjects_info,
		"data_rows":     job_dec.Data_rows,
		"last_update":   threeDaysAgoMidnight,
		"main_skill":    job_dec.Main_skill,
		"expand_skill":  job_dec.Expand_skill,
	}

	c.JSON(http.StatusOK, result)
}

func GetJobList(c *gin.Context) {
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

	var jobs []Job
	db_result := db.Find(&jobs)
	if db_result.Error != nil {
		log.Println("查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	var result []map[string]interface{}
	for _, job := range jobs {
		jobListMap := map[string]interface{}{
			"job_id":      job.ID,
			"job_type":    JOB_TYPE_MAP[job.Type],
			"job_type_id": job.Type,
			"job_name":    JOB_NAME_MAP[job.ID],
		}
		result = append(result, jobListMap)
	}

	c.JSON(http.StatusOK, result)
}

func GetRecomment(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("DealQ:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}
	var recomment_jobs struct {
		Job_arr string
		User_id int
	}

	db := baseClass.GetDB()
	tx := db.Begin()
	if db_result := tx.Table("recomment_jobs").Where("user_id = ?", userID).First(&recomment_jobs); db_result.Error == nil {
		strSlice := strings.Split(recomment_jobs.Job_arr, ",")

		// 将切片中的字符串转换为整数切片
		var intSlice []int
		for _, str := range strSlice {
			num, err := strconv.Atoi(str)
			if err != nil {
				log.Println("字符串转数组错误：", err)
			}
			intSlice = append(intSlice, num)
		}

		var result []map[string]interface{}
		for _, job := range intSlice {
			var job_info struct {
				Id   int
				Type int
			}
			tx.Table("jobs").Where("id = ?", job).First(&job_info)
			jobListMap := map[string]interface{}{
				"job_id":      job,
				"job_type":    JOB_TYPE_MAP[job_info.Type],
				"job_type_id": job_info.Type,
				"job_name":    JOB_NAME_MAP[job],
			}
			result = append(result, jobListMap)
		}
		tx.Commit()
		c.JSON(http.StatusOK, result)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"error": "无法找到用户数据"})
		return
	}
}
