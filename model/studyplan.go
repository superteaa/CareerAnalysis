package model

import (
	"CareerAnalysis/baseClass"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Study 模型
type Study struct {
	// ID         uint `gorm:"primaryKey"`
	UserID     int     `gorm:"column:user_id"`
	PlanName   string  `gorm:"column:plan_name" json:"plan_name"`
	SubjectID  int     `gorm:"column:subject_id" json:"subject_id"`
	StudyTime  int     `gorm:"column:study_time" json:"study_time"` //学习日期
	Spend_Time float64 `json:"spend_time"`                          // 学习时长
	AddTime    int     `gorm:"column:addtime" json:"add_time"`      // 用户操作的时间
	Note       string  `json:"note"`                                // 备注
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
	var xAxis []string
	date_info := make(map[string]interface{})
	subject_info := make(map[int]map[string]float64)

	// 获取当前时间并生成七天内的日期列表
	now := time.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := fmt.Sprintf("%d.%d", date.Month(), date.Day())
		xAxis = append(xAxis, dateStr)
		date_info[dateStr] = 1
	}

	// var date_info map[string]interface{}

	oldTime := now.AddDate(0, 0, -7)
	fiveAgo := time.Date(oldTime.Year(), oldTime.Month(), oldTime.Day(), 0, 0, 0, 0, oldTime.Location()).Unix()
	db.Where("user_id = ? and study_time >= ?", userID, fiveAgo).Find(&datas)

	var sum_time float64
	sum_time = 0

	// 计算所有 SubjectID 相同的 SpendTime 之和
	for _, v := range datas {
		// 将 StudyTime 转换为 time.Time 类型

		studyTime := int64(v.StudyTime)
		t := time.Unix(studyTime, 0)
		dateStr := fmt.Sprintf("%d.%d", t.Month(), t.Day())

		// 检查日期是否在 xAxis 列表中
		if _, exists := date_info[dateStr]; exists {
			// date_info[dateStr][v.SubjectID] += v.Spend_Time
			sum_time += v.Spend_Time
			if subject_info[v.SubjectID] == nil {
				subject_info[v.SubjectID] = make(map[string]float64)
			}
			subject_info[v.SubjectID][dateStr] += v.Spend_Time

		}
	}

	subjects_info := []map[string]interface{}{}
	for subjectID, v := range subject_info {
		dataArr := make([]string, 0, len(xAxis))
		for _, dateStr := range xAxis {
			if spendTime, exists := v[dateStr]; exists {
				dataDeal := fmt.Sprintf("%.1f", spendTime)
				dataArr = append(dataArr, dataDeal)
			} else {
				dataArr = append(dataArr, "0.0") // 若当天无数据，则填入0
			}
		}

		singleSubject := map[string]interface{}{
			"subject_id":   subjectID,
			"subject_name": SUBJECT_MAP[subjectID],
			"data":         dataArr,
		}

		subjects_info = append(subjects_info, singleSubject)
	}

	averageTime := fmt.Sprintf("%.2f", sum_time/7)

	result := map[string]interface{}{
		"subjects_info": subjects_info,
		"average_time":  averageTime,
		"xAxis":         xAxis,
	}

	c.JSON(http.StatusOK, result)
}
