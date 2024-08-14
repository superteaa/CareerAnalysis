package model

import (
	"CareerAnalysis/baseClass"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Study 模型
type Study struct {
	ID         int     `gorm:"primaryKey"`
	UserID     int     `gorm:"column:user_id"`
	PlanName   string  `gorm:"column:plan_name" json:"plan_name"`
	SubjectID  int     `gorm:"column:subject_id" json:"subject_id"`
	StudyTime  int     `gorm:"column:study_time" json:"study_time"` //学习日期
	Spend_Time float64 `json:"spend_time"`                          // 学习时长
	AddTime    int     `gorm:"column:addtime" json:"add_time"`      // 用户操作的时间
	Note       string  `json:"note"`                                // 备注
	Tags       string  `json:"tags"`                                // tag标签
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

var STUDY_TAG_MAP = map[int]string{
	1: "有问题",
	2: "新知识",
	3: "待总结",
	4: "没看懂",
	5: "有点问题",
	6: "其他",
}

func (Study) TableName() string {
	return "study_plans"
}

func AddPlan(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("AddPlan发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	_, ok := data["subject_id"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	_, ok = data["user_id"]
	if ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	_, ok = data["study_time"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	_, ok = data["spend_time"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	_, ok = data["addtime"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	_, ok = data["addtime"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("GetStudyList:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	data["user_id"] = userID

	tagsArr, ok := data["tags"].([]interface{})
	if ok {
		// tags 是一个 []interface{} 类型的数组
		isIntArray := true

		// 遍历数组元素，检查是否每个元素都是整数
		for _, item := range tagsArr {
			if _, ok := item.(float64); !ok {
				isIntArray = false
				log.Println(item)
				break
			}
		}

		if isIntArray {
			// tags 是一个整型数组

			// 将整数数组转换为字符串数组
			var strSlice []string
			for _, item := range tagsArr {
				num := item.(float64)
				strSlice = append(strSlice, strconv.FormatFloat(num, 'f', 0, 64))
			}

			// 使用逗号将字符串数组连接为一个字符串
			arrayStr := strings.Join(strSlice, ",")
			// arrayStr 现在包含了逗号分隔的整数数组字符串表示
			data["tags"] = arrayStr
		} else {
			// tags 不是整型数组
			log.Println("tags格式错误:")
			c.JSON(http.StatusOK, gin.H{"error": "tags格式错误"})
			return
		}
	}

	db := baseClass.GetDB()
	db_result := db.Table("study_plans").Create(&data)

	if db_result.Error != nil {
		log.Println("AddPlan数据库出错:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "success"})
	// if (data[])

}

func GetStudyData(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("GetStudyData发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	db := baseClass.GetDB()
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
		dataArr := make([]float64, 0, len(xAxis))
		for _, dateStr := range xAxis {
			if spendTime, exists := v[dateStr]; exists {

				dataArr = append(dataArr, spendTime)
			} else {
				dataArr = append(dataArr, 0.0) // 若当天无数据，则填入0
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

func GetPlanList(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("GetPlanList发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("GetTodayPlan:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	page := c.Query("page")
	pagesize := c.Query("pagesize")
	var pageInt int
	if page == "" {
		page = "1"
		pageInt, _ = strconv.Atoi(page)
	} else {
		pageInt, _ = strconv.Atoi(page)
	}
	var pagesizeInt int
	if pagesize == "" {
		pagesize = "10"
		pagesizeInt, _ = strconv.Atoi(pagesize)
	} else {
		pagesizeInt, _ = strconv.Atoi(pagesize)
	}

	var plan_list []Study
	// now := time.Now()
	// timestamp := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	db := baseClass.GetDB()

	db.Select("id, subject_id, spend_time, study_time, tags").Where("user_id = ?", userID).Order("study_time DESC").Order("addtime DESC").Limit(pagesizeInt).Offset((pageInt - 1) * pagesizeInt).Find(&plan_list)

	result := make(map[string][]map[string]interface{})

	for _, v := range plan_list {
		// 将字符串分割为切片
		strSlice := strings.Split(v.Tags, ",")

		// 将切片中的字符串转换为整数切片
		var intSlice []int

		for _, str := range strSlice {
			num, err := strconv.Atoi(str)
			if err != nil {
				log.Println("字符串转数组err：", err)
			}
			intSlice = append(intSlice, num)
		}

		// 将study_time转换为'YYYY-MM-DD'格式
		studyTimeStr := time.Unix(int64(v.StudyTime), 0).Format("2006-1-2")

		// 构建单条记录
		single := map[string]interface{}{
			"plan_id":    v.ID,
			"subject_id": v.SubjectID,
			"subject":    SUBJECT_MAP[v.SubjectID],
			"study_time": v.StudyTime,
			"spend_time": v.Spend_Time,
			"tags":       intSlice,
		}

		// 将记录添加到相应日期的切片中
		result[studyTimeStr] = append(result[studyTimeStr], single)
	}
	c.JSON(http.StatusOK, result)

}

func GetPlanDetail(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("GetPlanDetail发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("GetTodayPlan:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	plan_id := c.Query("plan_id")
	var planDetail Study
	db := baseClass.GetDB()
	if db_result := db.Where("id = ?", plan_id).First(&planDetail); db_result.Error != nil {
		log.Println("GetPlanDetail:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if planDetail.UserID != int(userID.(uint32)) {
		log.Println("GetPlanDetail:", "用户无权限访问该数据, planDetail:", planDetail.UserID, "userID:", userID)
		c.JSON(http.StatusOK, gin.H{"error": "用户无权限访问该数据"})
		return
	}

	// 将字符串分割为切片
	strSlice := strings.Split(planDetail.Tags, ",")

	// 将切片中的字符串转换为整数切片
	var intSlice []int

	for _, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			log.Println("字符串转数组err：", err)
		}
		intSlice = append(intSlice, num)
	}

	result := map[string]interface{}{
		"plan_id":    planDetail.ID,
		"plan_name":  planDetail.PlanName,
		"subject":    SUBJECT_MAP[planDetail.SubjectID],
		"subject_id": planDetail.SubjectID,
		"study_time": planDetail.StudyTime,
		"spend_time": planDetail.Spend_Time,
		"note":       planDetail.Note,
		"tags":       intSlice,
	}
	c.JSON(http.StatusOK, result)
}

func ChangePlan(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("GetPlanDetail发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("GetTodayPlan:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	plan_id := c.Query("plan_id")
	var planDetail Study
	db := baseClass.GetDB()
	if db_result := db.Where("id = ?", plan_id).First(&planDetail); db_result.Error != nil {
		log.Println("GetPlanDetail:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if planDetail.UserID != int(userID.(uint32)) {
		log.Println("GetPlanDetail:", "用户无权限访问该数据, planDetail:", planDetail.UserID, "userID:", userID)
		c.JSON(http.StatusOK, gin.H{"error": "用户无权限访问该数据"})
		return
	}

	// 将字符串分割为切片
	strSlice := strings.Split(planDetail.Tags, ",")

	// 将切片中的字符串转换为整数切片
	var intSlice []int

	for _, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			log.Println("字符串转数组err：", err)
		}
		intSlice = append(intSlice, num)
	}

	result := map[string]interface{}{
		"plan_id":    planDetail.ID,
		"plan_name":  planDetail.PlanName,
		"subject":    SUBJECT_MAP[planDetail.SubjectID],
		"subject_id": planDetail.SubjectID,
		"study_time": planDetail.StudyTime,
		"spend_time": planDetail.Spend_Time,
		"note":       planDetail.Note,
		"tags":       intSlice,
	}
	c.JSON(http.StatusOK, result)
}
