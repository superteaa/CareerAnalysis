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
	SubjectCatKey int     `gorm:"column:subject_cat_key"`
	SubjectSubKey int     `gorm:"column:subject_sub_key"`
	SubjectKey    int     `gorm:"column:subject_key"`
	StudyTime     int     `gorm:"column:study_time" json:"study_time"` //学习日期
	Spend_Time    float64 `json:"spend_time"`                          // 学习时长
	AddTime       int     `gorm:"column:addtime" json:"add_time"`      // 用户操作的时间
	Note          string  `json:"note"`                                // 备注
	Tags          string  `json:"tags"`                                // tag标签
}

var SUBJECT_MAP = map[int]map[int]map[int]string{
	1: { // 软件类
		1: { // 编程语言
			1: "Java",
			2: "Python",
			3: "C++",
			4: "JavaScript",
			5: "Ruby",
			6: "PHP",
			7: "Go",
		},
		2: { // 前端开发
			1: "HTML",
			2: "CSS",
			3: "JavaScript",
			4: "Vue.js",
			5: "React.js",
			6: "webpack",
			7: "vite",
			8: "router",
			9: "Bootstrap",
		},
		3: { // 后端开发框架
			1: "Spring",
			2: "Django",
			3: "Express",
			4: "Yii",
			5: "Gin",
		},
		4: { // 数据库
			1: "MySQL",
			2: "PostgreSQL",
			3: "MongoDB",
			4: "Redis",
			5: "Oracle",
			6: "SQLite",
			7: "SQL Server",
		},
		5: { // 容器技术
			1: "Docker",
			2: "Kubernetes",
			3: "Jenkins",
			4: "Terraform",
		},
		6: { // 云计算
			1: "AWS",
			2: "Azure",
			3: "Google Cloud",
			4: "OpenStack",
			5: "Alibaba Cloud",
		},
		7: { // 测试工具
			1: "Selenium",
			2: "TestNG",
			3: "Postman",
			4: "Cucumber",
		},
		8: { // 版本控制
			1: "Git",
			2: "SVN",
			3: "Mercurial",
		},
		9: { // 操作系统
			1: "Linux",
			2: "Windows",
			3: "macOS",
		},
		10: { // 移动开发
			1: "Android开发（Kotlin/Java）",
			2: "iOS开发（Swift/Objective-C）",
			3: "React Native",
			4: "Flutter",
		},
	},
	2: { // 硬件类
		1: { // 嵌入式系统
			1: "ARM",
			2: "Raspberry Pi",
			3: "Arduino",
			4: "FPGA",
			5: "STM32",
			6: "ESP32",
		},
		2: { // 硬件设计工具
			1: "Verilog",
			2: "VHDL",
			3: "Altium Designer",
			4: "OrCAD",
			5: "KiCad",
			6: "Proteus",
		},
		3: { // 微处理器
			1: "Intel x86",
			2: "ARM Cortex",
			3: "PIC",
			4: "AVR",
			5: "MSP430",
		},
		4: { // 通信接口
			1: "UART",
			2: "I2C",
			3: "SPI",
			4: "CAN",
			5: "USB",
			6: "Ethernet",
		},
		5: { // 传感器技术
			1: "ADC/DAC",
			2: "温度传感器",
			3: "压力传感器",
			4: "加速度计",
			5: "陀螺仪",
		},
		6: { // 电源管理
			1: "电池管理",
			2: "电源调节",
			3: "充电管理IC",
		},
		7: { // 嵌入式软件开发工具
			1: "Keil",
			2: "IAR Embedded Workbench",
			3: "Atmel Studio",
			4: "MPLAB",
		},
		8: { // 调试工具
			1: "逻辑分析仪",
			2: "示波器",
			3: "JTAG调试器",
			4: "ICE仿真器",
		},
		9: { // 生产工艺
			1: "SMT",
			2: "焊接技术",
			3: "自动光学检测(AOI)",
		},
		10: { // 驱动程序开发
			1: "驱动程序开发（Linux Kernel）",
			2: "驱动程序开发（Windows Driver Model）",
		},
	},
	3: { // 网络类
		1: { // 网络协议
			1: "TCP/IP",
			2: "HTTP/HTTPS",
			3: "DNS",
			4: "DHCP",
			5: "BGP",
			6: "OSPF",
			7: "MPLS",
		},
		2: { // 网络设备
			1: "路由器",
			2: "交换机",
			3: "防火墙",
			4: "负载均衡器",
			5: "VPN设备",
		},
		3: { // 网络监控工具
			1: "Wireshark",
			2: "Nagios",
			3: "Zabbix",
			4: "SolarWinds",
			5: "Cacti",
			6: "NetFlow",
		},
		4: { // 网络虚拟化
			1: "SDN（软件定义网络）",
			2: "NFV（网络功能虚拟化）",
			3: "VLAN",
			4: "VXLAN",
		},
		5: { // 网络安全
			1: "防火墙",
			2: "IDS/IPS",
			3: "WAF（Web应用防火墙）",
			4: "DDoS防护",
			5: "VPN",
			6: "SSL/TLS",
			7: "加密技术",
			8: "SOC（安全运营中心）",
			9: "SIEM（安全信息和事件管理）",
		},
		6: { // 无线通信技术
			1: "Wi-Fi",
			2: "Bluetooth",
			3: "ZigBee",
			4: "LoRa",
			5: "5G",
		},
		7: { // 网络操作系统
			1: "Cisco IOS",
			2: "Juniper Junos",
			3: "Huawei VRP",
			4: "Ansible（网络自动化）",
		},
		8: { // 云网络
			1: "云端虚拟私有网络（VPC）",
			2: "负载均衡",
			3: "DNS服务（AWS Route 53，Azure DNS）",
		},
	},
	4: { // 信息系统类
		1: { // ERP系统
			1: "SAP ERP",
			2: "Oracle E-Business Suite",
			3: "Microsoft Dynamics",
			4: "Odoo",
			5: "金蝶",
			6: "用友",
		},
		2: { // CRM系统
			1: "Salesforce",
			2: "HubSpot",
			3: "Zoho CRM",
			4: "SugarCRM",
			5: "Dynamics CRM",
		},
		3: { // 数据仓库
			1: "Amazon Redshift",
			2: "Google BigQuery",
			3: "Snowflake",
			4: "Apache Hive",
			5: "Teradata",
		},
		4: { // 数据集成工具
			1: "Apache Nifi",
			2: "Talend",
			3: "Informatica",
			4: "Microsoft SSIS",
			5: "Pentaho",
		},
		5: { // 商业智能
			1: "Tableau",
			2: "Power BI",
			3: "Looker",
			4: "QlikView",
			5: "Sisense",
		},
		6: { // 消息中间件
			1: "Apache Kafka",
			2: "RabbitMQ",
			3: "ActiveMQ",
			4: "IBM MQ",
			5: "Tibco",
		},
		7: { // 身份管理
			1: "LDAP",
			2: "Active Directory",
			3: "Okta",
			4: "Auth0",
			5: "AWS IAM",
		},
		8: { // 大数据处理
			1: "Apache Hadoop",
			2: "Apache Spark",
			3: "Flink",
			4: "Presto",
			5: "HDFS",
			6: "NoSQL数据库",
		},
		9: { // 企业应用集成
			2: "MuleSoft",
			3: "WSO2",
			4: "Oracle SOA Suite",
			5: "TIBCO BusinessWorks",
		},
		10: { // 内容管理系统
			1: "WordPress",
			2: "Joomla!",
			3: "Drupal",
			4: "Wix",
			5: "Squarespace",
		},
	},
	5: { // 制造类
		1: { // 自动化控制
			1: "PLC编程",
			2: "SCADA系统",
			3: "DCS系统",
		},
		2: { // 工业机器人
			1: "机械臂编程",
			2: "工业物联网",
			3: "ROS",
		},
		3: { // 嵌入式系统
			1: "实时操作系统",
			2: "工业控制单片机",
			3: "嵌入式Linux",
			4: "FreeRTOS",
		},
		4: { // CAD/CAM软件
			1: "AutoCAD",
			2: "SolidWorks",
			3: "Fusion 360",
			4: "CATIA",
			5: "Siemens NX",
		},
		5: { // 工业物联网
			1: "OPC UA",
			2: "MQTT",
			3: "工业以太网",
			4: "Profinet",
			5: "Modbus",
		},
	},
}

// 定义分类名称映射
var CATEGORY_NAMES = map[int]string{
	1: "软件类",
	2: "硬件类",
	3: "网络类",
	4: "信息系统类",
	5: "制造类",
}

// 定义子分类名称映射
var SUBCATEGORY_NAMES = map[int]map[int]string{
	1: { // 软件类
		1:  "编程语言",
		2:  "前端开发",
		3:  "后端开发框架",
		4:  "数据库",
		5:  "容器技术",
		6:  "云计算",
		7:  "测试工具",
		8:  "版本控制",
		9:  "操作系统",
		10: "移动开发",
	},
	2: { // 硬件类
		1:  "嵌入式系统",
		2:  "硬件设计工具",
		3:  "微处理器",
		4:  "通信接口",
		5:  "传感器技术",
		6:  "电源管理",
		7:  "嵌入式软件开发工具",
		8:  "调试工具",
		9:  "生产工艺",
		10: "驱动程序开发",
	},
	3: { // 网络类
		1: "网络协议",
		2: "网络设备",
		3: "网络监控工具",
		4: "网络虚拟化",
		5: "网络安全",
		6: "无线通信技术",
		7: "网络操作系统",
		8: "云网络",
	},
	4: { // 信息系统类
		1:  "ERP系统",
		2:  "CRM系统",
		3:  "数据仓库",
		4:  "数据集成工具",
		5:  "商业智能",
		6:  "消息中间件",
		7:  "身份管理",
		8:  "大数据处理",
		9:  "企业应用集成",
		10: "内容管理系统",
	},
	5: { // 制造类
		1: "自动化控制",
		2: "工业机器人",
		3: "嵌入式系统",
		4: "CAD/CAM软件",
		5: "工业物联网",
	},
}

func mapToID(category, subCategory, value string) (int, int, int) {
	for catKey, subMap := range SUBJECT_MAP {
		if CATEGORY_NAMES[catKey] == category {
			for subKey, valueMap := range subMap {
				if SUBCATEGORY_NAMES[catKey][subKey] == subCategory {
					for valueKey, valueName := range valueMap {
						if valueName == value {
							return catKey, subKey, valueKey
						}
					}
				}
			}
		}
	}
	return 0, 0, 0
}

// func main() {
// 	// 示例用法
// 	category := "软件类"
// 	subCategory := "编程语言"
// 	value := "Go"

// 	catID, subID, valueID, err := mapToID(category, subCategory, value)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Printf("Category ID: %d, SubCategory ID: %d, Value ID: %d\n", catID, subID, valueID)
// 	}
// }

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
	PlanName      string        `json:"plan_name"`
	SubjectCatKey string        `json:"subject_cat_key" binding:"required"`
	SubjectSubKey string        `json:"subject_sub_key" binding:"required"`
	SubjectKey    string        `json:"subject_key" binding:"required"`
	StudyTime     int           `json:"study_time" binding:"required"`
	SpendTime     float64       `json:"spend_time" binding:"required"`
	AddTime       int           `json:"add_time" binding:"required"`
	Note          string        `json:"note"`
	Tags          []interface{} `json:"tags"`
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

	// 将分类、子分类和键值转换为对应的ID
	catKey, subKey, key := mapToID(req.SubjectCatKey, req.SubjectSubKey, req.SubjectKey)
	if catKey == 0 && subKey == 0 && key == 0 {
		log.Println("无效的科目信息", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的科目信息"})
		return
	}

	// 构造数据
	study := Study{
		UserID:        int(userID.(uint32)),
		PlanName:      req.PlanName,
		SubjectCatKey: catKey,
		SubjectSubKey: subKey,
		SubjectKey:    key,
		StudyTime:     req.StudyTime,
		Spend_Time:    req.SpendTime,
		AddTime:       req.AddTime,
		Note:          req.Note,
		Tags:          tagsStr,
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

			// 获取学科名称
			subjectName := SUBJECT_MAP[v.SubjectCatKey][v.SubjectSubKey][v.SubjectKey]
			if subject_info[subjectName] == nil {
				subject_info[subjectName] = make(map[string]float64)
			}
			subject_info[subjectName][dateStr] += v.Spend_Time
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

		// 获取学科名称
		subjectName := SUBJECT_MAP[v.SubjectCatKey][v.SubjectSubKey][v.SubjectKey]

		// 将study_time转换为'YYYY-MM-DD'格式
		studyTimeStr := time.Unix(int64(v.StudyTime), 0).Format("2006-1-2")

		// 构建单条记录
		single := map[string]interface{}{
			"plan_id":    v.ID,
			"subject":    subjectName,
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

	// 获取学科名称
	// subjectName := SUBJECT_MAP[planDetail.SubjectCatKey][planDetail.SubjectSubKey][planDetail.SubjectKey]

	result := map[string]interface{}{
		"plan_id":         planDetail.ID,
		"plan_name":       planDetail.PlanName,
		"subject_cat_key": CATEGORY_NAMES[planDetail.SubjectCatKey],
		"subject_sub_key": SUBCATEGORY_NAMES[planDetail.SubjectCatKey][planDetail.SubjectSubKey],
		"subject_key":     SUBJECT_MAP[planDetail.SubjectCatKey][planDetail.SubjectSubKey][planDetail.SubjectKey],
		"study_time":      planDetail.StudyTime,
		"spend_time":      planDetail.Spend_Time,
		"note":            planDetail.Note,
		"tags":            intSlice,
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

	requiredFields := []string{"plan_id", "study_time", "spend_time", "add_time", "subject_cat_key", "subject_sub_key", "subject_key"}
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
	planDetail.SubjectCatKey, planDetail.SubjectSubKey, planDetail.SubjectKey = mapToID(data["subject_cat_key"].(string), data["subject_sub_key"].(string), data["subject_key"].(string))
	// planDetail.SubjectCatKey = int(data["subject_cat_key"].(float64))
	// planDetail.SubjectSubKey = int(data["subject_sub_key"].(float64))
	// planDetail.SubjectKey = int(data["subject_key"].(float64))

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

	result := make(map[string]map[string][]string)

	for categoryID, subcategories := range SUBJECT_MAP {
		categoryName, categoryExists := CATEGORY_NAMES[categoryID]
		if !categoryExists {
			continue
		}

		subcategoryMap := make(map[string][]string)
		for subcategoryID, subjects := range subcategories {
			subcategoryName, subcategoryExists := SUBCATEGORY_NAMES[categoryID][subcategoryID]
			if !subcategoryExists {
				continue
			}

			var subjectMap []string
			for _, subjectName := range subjects {
				subjectMap = append(subjectMap, subjectName)
			}

			subcategoryMap[subcategoryName] = subjectMap
		}

		result[categoryName] = subcategoryMap
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
		SubjectCatKey int     `gorm:"column:subject_cat_key" json:"subject_cat_key"`
		SubjectSubKey int     `gorm:"column:subject_sub_key" json:"subject_sub_key"`
		SubjectKey    int     `gorm:"column:subject_key" json:"subject_key"`
		SpendTime     float64 `gorm:"column:spend_time" json:"spend_time"`
	}

	// 查询有学习记录的科目
	var learnedSkills []LearnedSkill
	db := baseClass.GetDB()
	db.Raw("SELECT subject_cat_key, subject_sub_key, subject_key, spend_time FROM study_plans WHERE user_id = ? AND spend_time > 0", userID).Scan(&learnedSkills)

	// 组织返回数据
	result := make(map[string]map[string][]string)

	for _, skill := range learnedSkills {
		// 获取分类名
		categoryName, categoryExists := CATEGORY_NAMES[skill.SubjectCatKey]
		if !categoryExists {
			continue
		}

		// 获取子分类名
		subcategoryName, subcategoryExists := SUBCATEGORY_NAMES[skill.SubjectCatKey][skill.SubjectSubKey]
		if !subcategoryExists {
			continue
		}

		// 获取科目名称
		subjectName, subjectExists := SUBJECT_MAP[skill.SubjectCatKey][skill.SubjectSubKey][skill.SubjectKey]
		if !subjectExists {
			continue
		}

		// 如果分类不存在则初始化
		if _, exists := result[categoryName]; !exists {
			result[categoryName] = make(map[string][]string)
		}

		// 如果子分类不存在则初始化
		if _, exists := result[categoryName][subcategoryName]; !exists {
			result[categoryName][subcategoryName] = []string{}
		}

		// 检查是否已经存在该科目名称
		if !contains(result[categoryName][subcategoryName], subjectName) {
			result[categoryName][subcategoryName] = append(result[categoryName][subcategoryName], subjectName)
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
