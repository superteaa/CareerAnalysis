package model

import (
	"CareerAnalysis/baseClass"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Skill 模型
type Skill struct {
	// ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"column:user_id"`
	Subject    string
	Spend_Time uint
	// Date       string
}

func GetSkillList(c *gin.Context) {
	db := baseClass.InitDB()
	userID, exists := c.Get("userID")

	if !exists {
		log.Println("GetSkillList:", "用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	var skills []Skill

	var subjects_info []map[string]interface{}

	db.Where("user_id = ?", userID).Find(&skills)

	var sum_time uint
	for _, v := range skills {
		sum_time = sum_time + v.Spend_Time
		subject_info := map[string]interface{}{
			"subject_name":  v.Subject,
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
