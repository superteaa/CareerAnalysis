package model

import (
	"CareerAnalysis/baseClass"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Study 模型
type Study struct {
	// ID         uint `gorm:"primaryKey"`
	UserID     int    `gorm:"column:user_id"`
	PlanName   string `gorm:"column:plan_name"`
	SubjectID  int    `gorm:"column:subject_id"`
	Spend_Time int
	AddTime    int
}

func (Study) TableName() string {
	return "study_plans"
}

func AddPlan(c *gin.Context) {

}

func GetStudyList(c *gin.Context) {
	db := baseClass.InitDB()
	userID, exists := c.Get("userID")

	if !exists {
		log.Println("GetStudyList:", "用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	var skills []Study

	var subjects_info []map[string]interface{}

	db.Where("user_id = ?", userID).Find(&skills)

	var sum_time int
	for _, v := range skills {
		sum_time = sum_time + v.Spend_Time
		subject_info := map[string]interface{}{
			"subject_name":  v.SubjectID,
			"subject_spend": v.Spend_Time,
		}
		subjects_info = append(subjects_info, subject_info)
	}

	result := map[string]interface{}{
		"subjects_info": subjects_info,
		"sum_time":      sum_time,
	}

	c.JSON(http.StatusOK, result)
}
