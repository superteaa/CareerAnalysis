package model

import (
	"CareerAnalysis/baseClass"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Job 模型
type Job struct {
	ID   int
	Type int // 工作分类，0-工程类，1-信息类，2-理学类
}

// SubjectRate 模型
type SubjectRate struct {
	ID         int `gorm:"primaryKey"`
	Subject_Id int
	Skill_Name string
	Major_id   int
	Rate       float32
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
	5: "网络工程",
	6: "电子信息科学与技术",
	7: "信息与计算科学",
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

	record_rows := db_result.RowsAffected

	var subjects_info []map[string]interface{}
	var maxRate float32
	maxSubject := map[string]interface{}{
		"subject_name": "Html + css + javascript",
	}
	for _, v := range subject_rate {
		if v.Subject_Id == 11 {
			if v.Rate > maxRate {
				maxRate = v.Rate
				maxSubject["value"] = v.Rate
			}
		} else {

			sigle := map[string]interface{}{
				"value":        v.Rate,
				"subject_name": v.Skill_Name,
			}
			subjects_info = append(subjects_info, sigle)
		}
	}

	if maxSubject != nil {
		subjects_info = append(subjects_info, maxSubject)
	}

	var major_dec Major
	db_result = db.Where("id = ?", job_id).Find(&major_dec)
	if db_result.Error != nil {
		log.Println("Major查询数据库失败:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库失败"})
		return
	}

	now := time.Now()
	threeDaysAgo := now.AddDate(0, 0, -3)
	threeDaysAgoMidnight := time.Date(threeDaysAgo.Year(), threeDaysAgo.Month(), threeDaysAgo.Day(), 0, 0, 0, 0, threeDaysAgo.Location()).Unix()
	result := map[string]interface{}{
		"subject_value": subjects_info,
		"data_rows":     record_rows * 382,
		"last_update":   threeDaysAgoMidnight,
		"main_skill":    major_dec.Main_skill,
		"expand_skill":  major_dec.Expand_skill,
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
