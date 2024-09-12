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
	ID       int    `gorm:"primaryKey"`
	UserID   int    `gorm:"column:user_id"`
	PlanName string `gorm:"column:plan_name" json:"plan_name"`
	// SubjectID  int     `gorm:"column:subject_id" json:"subject_id"`
	JobId         int     `gorm:"column:job_id"`
	SkillId       int     `gorm:"column:skill_id"`
	SubjectCatKey int     `gorm:"column:subject_cat_key"`
	SubjectSubKey int     `gorm:"column:subject_sub_key"`
	SubjectKey    int     `gorm:"column:subject_key"`
	StudyTime     int     `gorm:"column:study_time" json:"study_time"` //学习日期
	Spend_Time    float64 `json:"spend_time"`                          // 学习时长
	AddTime       int     `gorm:"column:addtime" json:"add_time"`      // 用户操作的时间
	Note          string  `json:"note"`                                // 备注
	Tags          string  `json:"tags"`                                // tag标签
}

type Skill struct {
	JobID     int    `gorm:"column:job_id"`
	ID        int    `gorm:"column:id"`
	SkillName string `gorm:"column:skill_name"`
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

type AddPlanRequest struct {
	PlanName string `json:"plan_name"`
	JobId    int    `json:"job_id" binding:"required"`
	SkillId  int    `json:"skill_id" binding:"required"`
	// SubjectKey string        `json:"subject_key" binding:"required"`
	StudyTime int           `json:"study_time" binding:"required"`
	SpendTime float64       `json:"spend_time" binding:"required"`
	AddTime   int           `json:"add_time" binding:"required"`
	Note      string        `json:"note"`
	Tags      []interface{} `json:"tags"`
}

// 处理标签，将标签数组转换为逗号分隔的字符串
func processTags(tags []interface{}) (string, error) {
	if tags == nil {
		return "", nil
	}

	var strSlice []string
	for _, item := range tags {
		if num, ok := item.(float64); ok {
			strSlice = append(strSlice, strconv.FormatFloat(num, 'f', 0, 64))
		} else {
			return "", fmt.Errorf("标签必须为数字类型")
		}
	}

	return strings.Join(strSlice, ","), nil
}

func AddPlan(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("AddPlan发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	// 解析并验证请求参数
	var req AddPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("AddPlan请求参数错误:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("AddPlan: 鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	// 处理标签
	tagsStr, err := processTags(req.Tags)
	if err != nil {
		log.Println("标签格式错误:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "tags格式错误"})
		return
	}

	// // 将分类、子分类和键值转换为对应的ID
	// catKey, subKey, key := mapToID(req.SubjectCatKey, req.SubjectSubKey, req.SubjectKey)
	// if catKey == 0 && subKey == 0 && key == 0 {
	// 	log.Println("无效的科目信息", userID)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "无效的科目信息"})
	// 	return
	// }

	// 构造数据
	study := Study{
		UserID:   int(userID.(uint32)),
		PlanName: req.PlanName,
		// SubjectCatKey: catKey,
		// SubjectSubKey: subKey,
		// SubjectKey:    key,
		JobId:      req.JobId,
		SkillId:    req.SkillId,
		StudyTime:  req.StudyTime,
		Spend_Time: req.SpendTime,
		AddTime:    req.AddTime,
		Note:       req.Note,
		Tags:       tagsStr,
	}

	// 启动事务
	db := baseClass.GetDB()
	tx := db.Begin()

	// 插入数据
	if err := tx.Create(&study).Error; err != nil {
		log.Println("AddPlan数据库出错:", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
		return
	}

	// 提交事务
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
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
	subject_info := make(map[string]map[string]float64)

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

	// 计算所有相同 Subject 的 SpendTime 之和
	for _, v := range datas {
		studyTime := int64(v.StudyTime)
		t := time.Unix(studyTime, 0)
		dateStr := fmt.Sprintf("%d.%d", t.Month(), t.Day())

		if _, exists := date_info[dateStr]; exists {
			sum_time += v.Spend_Time

			// 获取技能名
			var find_skill Skill
			db = baseClass.GetDB()

			// 保存更改
			if err := db.Table("skills").Where("id = ?", v.SkillId).First(&find_skill).Error; err != nil {
				log.Println("GetSkillTree数据库错误:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
				return
			}

			if subject_info[find_skill.SkillName] == nil {
				subject_info[find_skill.SkillName] = make(map[string]float64)
			}
			subject_info[find_skill.SkillName][dateStr] += v.Spend_Time
		}
	}

	subjects_info := []map[string]interface{}{}
	for subjectName, v := range subject_info {
		dataArr := make([]float64, 0, len(xAxis))
		for _, dateStr := range xAxis {
			if spendTime, exists := v[dateStr]; exists {
				dataArr = append(dataArr, spendTime)
			} else {
				dataArr = append(dataArr, 0.0) // 若当天无数据，则填入0
			}
		}

		singleSubject := map[string]interface{}{
			"subject_name": subjectName,
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
	db := baseClass.GetDB()

	db.Select("id, subject_cat_key, subject_sub_key, subject_key, spend_time, study_time, tags").
		Where("user_id = ?", userID).
		Order("study_time DESC").
		Order("addtime DESC").
		Limit(pagesizeInt).
		Offset((pageInt - 1) * pagesizeInt).
		Find(&plan_list)

	result := make(map[string][]map[string]interface{})

	for _, v := range plan_list {
		// 将字符串分割为切片
		strSlice := strings.Split(v.Tags, ",")

		// 将切片中的字符串转换为整数切片
		var intSlice []int
		for _, str := range strSlice {
			num, err := strconv.Atoi(str)
			if err != nil {
				log.Println("字符串转数组错误：", err)
			}
			intSlice = append(intSlice, num)
		}

		// 获取技能名
		var find_skill Skill
		db = baseClass.GetDB()

		// 保存更改
		if err := db.Table("skills").Where("id = ?", v.SkillId).First(&find_skill).Error; err != nil {
			log.Println("GetSkillTree数据库错误:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
			return
		}

		// 将study_time转换为'YYYY-MM-DD'格式
		studyTimeStr := time.Unix(int64(v.StudyTime), 0).Format("2006-1-2")

		// 构建单条记录
		single := map[string]interface{}{
			"plan_id":    v.ID,
			"subject":    find_skill.SkillName,
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
		log.Println("GetPlanDetail:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	plan_id := c.Query("plan_id")
	if plan_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}
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
			log.Println("字符串转数组错误：", err)
		}
		intSlice = append(intSlice, num)
	}

	// 获取技能名
	var find_skill Skill
	db = baseClass.GetDB()

	// 保存更改
	if err := db.Table("skills").Where("id = ?", planDetail.SkillId).First(&find_skill).Error; err != nil {
		log.Println("GetSkillTree数据库错误:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
		return
	}

	result := map[string]interface{}{
		"plan_id":    planDetail.ID,
		"plan_name":  planDetail.PlanName,
		"job_name":   JOB_NAME_MAP[planDetail.JobId],
		"skill_name": find_skill.SkillName,
		// "subject_key":     SUBJECT_MAP[planDetail.SubjectCatKey][planDetail.SubjectSubKey][planDetail.SubjectKey],
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
			log.Println("ChangePlan发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		log.Println("参数绑定错误:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	requiredFields := []string{"plan_id", "study_time", "spend_time", "add_time", "job_id", "skill_id"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数"})
			return
		}
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("ChangePlan:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	// 处理 tags
	if tagsArr, ok := data["tags"].([]interface{}); ok {
		var strSlice []string
		for _, item := range tagsArr {
			if num, ok := item.(float64); ok {
				strSlice = append(strSlice, strconv.FormatFloat(num, 'f', 0, 64))
			} else {
				log.Println("tags格式错误:", item)
				c.JSON(http.StatusBadRequest, gin.H{"error": "tags格式错误"})
				return
			}
		}
		data["tags"] = strings.Join(strSlice, ",")
	}

	// 查询计划详情
	var planDetail Study
	db := baseClass.GetDB()
	if err := db.Where("id = ?", data["plan_id"]).First(&planDetail).Error; err != nil {
		log.Println("ChangePlan数据库查询错误:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
		return
	}

	// 检查用户权限
	if planDetail.UserID != int(userID.(uint32)) {
		log.Println("ChangePlan:", "用户无权限访问该数据, planDetail:", planDetail.UserID, "userID:", userID)
		c.JSON(http.StatusOK, gin.H{"error": "用户无权限访问该数据"})
		return
	}

	// 更新计划详情
	if val, ok := data["plan_name"].(string); ok {
		planDetail.PlanName = val
	}
	planDetail.JobId, planDetail.SkillId = int(data["job_id"].(float64)), int(data["skill_id"].(float64))

	planDetail.StudyTime = int(data["study_time"].(float64))
	planDetail.Spend_Time = data["spend_time"].(float64)
	planDetail.AddTime = int(data["add_time"].(float64))

	if val, ok := data["note"].(string); ok {
		planDetail.Note = val
	}
	if val, ok := data["tags"].(string); ok {
		planDetail.Tags = val
	}

	// 保存更改
	if err := db.Save(&planDetail).Error; err != nil {
		log.Println("ChangePlan数据库保存错误:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

// 将 SUBJECT_MAP 转换为带有分类名称的格式
func GetSubjectMap(c *gin.Context) {

	var skills []Skill
	db := baseClass.GetDB()
	if err := db.Find(&skills).Error; err != nil {
		log.Println("GetSubjectMap数据库错误:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
		return
	}

	SUBJECT_MAP := make(map[int]map[int]string)

	for _, skill := range skills {
		if _, ok := SUBJECT_MAP[skill.JobID]; !ok {
			SUBJECT_MAP[skill.JobID] = make(map[int]string)
		}
		SUBJECT_MAP[skill.JobID][skill.ID] = skill.SkillName
	}

	result := make(map[int]map[string]interface{})

	for jobID, skills := range SUBJECT_MAP {
		jobName, jobExists := JOB_NAME_MAP[jobID]
		if !jobExists {
			continue
		}
		tmp := map[string]interface{}{
			"job_name": jobName,
			"skills":   skills,
		}
		result[jobID] = tmp
	}

	c.JSON(http.StatusOK, result)
}

func GetSkillTree(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("GetSkillTree: 鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

	type LearnedSkill struct {
		JobId     int     `gorm:"column:job_id" json:"job_id"`
		SkillId   int     `gorm:"column:skill_id" json:"skill_id"`
		SpendTime float64 `gorm:"column:spend_time" json:"spend_time"`
	}

	// 查询有学习记录的科目
	var learnedSkills []LearnedSkill
	db := baseClass.GetDB()
	db.Raw("SELECT job_id, skill_id, spend_time FROM study_plans WHERE user_id = ? AND spend_time > 0", userID).Scan(&learnedSkills)

	// 组织返回数据
	result := make(map[int]map[string]interface{})

	for _, skill := range learnedSkills {
		// 获取岗位名
		jobName, jobExists := JOB_NAME_MAP[skill.JobId]
		if !jobExists {
			continue
		}

		// 获取技能名
		var find_skill Skill
		db := baseClass.GetDB()

		// 保存更改
		if err := db.Table("skills").Where("id = ?", skill.SkillId).First(&find_skill).Error; err != nil {
			log.Println("GetSkillTree数据库错误:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
			return
		}

		// 如果分类不存在则初始化
		if _, exists := result[skill.JobId]; !exists {
			result[skill.JobId] = make(map[string]interface{})
			result[skill.JobId]["job_name"] = jobName
			result[skill.JobId]["skills"] = []string{} // 初始化 skills 列表
		}

		// 检查是否已经存在该科目名称
		if skills, ok := result[skill.JobId]["skills"].([]string); ok {
			if !contains(skills, find_skill.SkillName) {
				result[skill.JobId]["skills"] = append(skills, find_skill.SkillName)
			}
		}

	}

	c.JSON(http.StatusOK, result)
}

// 辅助函数，用于检查切片中是否包含特定元素
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
