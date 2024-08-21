package model

import (
	"CareerAnalysis/baseClass"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SubjectRate 模型
type SubjectRate struct {
	ID         int `gorm:"primaryKey"`
	Subject_Id int
	Skill_Name string
	Major_id   int
	Rate       float32
}

func (SubjectRate) TableName() string {
	return "skills"
}

func GetSubjectRate(c *gin.Context) {
	log.Println("v.ID")
	major_id := c.Query("major_id")

	db := baseClass.GetDB()

	var subject_rate []SubjectRate
	db_result := db.Where("major_id = ?", major_id).Find(&subject_rate)
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
	db_result = db.Where("id = ?", major_id).Find(&major_dec)
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
