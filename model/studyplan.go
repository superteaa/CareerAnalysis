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
	PlanName   string `gorm:"column:plan_name" json:"plan_name"`
	SubjectID  int    `gorm:"column:subject_id" json:"subject_id"`
	StudyTime  int    `gorm:"column:study_time" json:"study_time"` //学习日期
	Spend_Time int    `json:"spend_time"`                          // 学习时长
	AddTime    int    `gorm:"column:addtime" json:"add_time"`      // 用户操作的时间
	Note       string `json:"note"`                                // 备注
}

var SUBJECT_MAP = map[int]string{
	1: "Java",
	2: "C语言",
	3: "Python",
	4: "C++",
	// 5: "网络工程",
	// 6: "电子信息科学与技术",
	// 7: "信息与计算科学",
}

func (Study) TableName() string {
	return "study_plans"
}

func AddPlan(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	var data Study
	if err := c.BindJSON(&data); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}
	if data.SubjectID == 0 || data.StudyTime == 0 || data.Spend_Time == 0 || data.AddTime == 0 {
		log.Println(data)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("GetStudyList:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	data.UserID = userID.(int)

	db := baseClass.InitDB()
	db.Create(&data)

	c.JSON(http.StatusOK, gin.H{"msg": "success"})
	// if (data[])

}

func GetStudyData(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	db := baseClass.InitDB()
	userID, exists := c.Get("userID")

	if !exists {
		log.Println("GetStudyList:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	var datas []Study

	var subjects_info []map[string]interface{}

	db.Where("user_id = ?", userID).Find(&datas)

	sum_time := 0
	subjectSpendMap := make(map[int]int) // 用于存储每个 SubjectID 的 SpendTime 之和

	// 计算所有 SubjectID 相同的 SpendTime 之和
	for _, v := range datas {
		sum_time += v.Spend_Time
		subjectSpendMap[v.SubjectID] += v.Spend_Time
	}

	// 构建 subjects_info 列表
	for subjectID, spendTime := range subjectSpendMap {
		subject_info := map[string]interface{}{
			"subject_name":  SUBJECT_MAP[subjectID],
			"subject_id":    subjectID,
			"subject_spend": spendTime,
		}
		subjects_info = append(subjects_info, subject_info)
	}

	result := map[string]interface{}{
		"subjects_info": subjects_info,
		"sum_time":      sum_time,
	}

	c.JSON(http.StatusOK, result)
}
