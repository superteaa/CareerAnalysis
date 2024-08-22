package model

import (
	"CareerAnalysis/baseClass"
	"errors"
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
	// SubjectID  int     `gorm:"column:subject_id" json:"subject_id"`
	SubjectCatKey  int     `gorm:"column:subject_cat_key"`
	SubjectSubKey  int     `gorm:"column:subject_sub_key"`
	SubjectKey  int     `gorm:"column:subject_key"`
	StudyTime  int     `gorm:"column:study_time" json:"study_time"` //学习日期
	Spend_Time float64 `json:"spend_time"`                          // 学习时长
	AddTime    int     `gorm:"column:addtime" json:"add_time"`      // 用户操作的时间
	Note       string  `json:"note"`                                // 备注
	Tags       string  `json:"tags"`                                // tag标签
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

func mapToID(category, subCategory, value string) (int, int, int) {
	for catKey, subMap := range SUBJECT_MAP {
		if catKeyName(catKey) == category {
			for subKey, valueMap := range subMap {
				if subKeyName(catKey, subKey) == subCategory {
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

// 辅助函数，用于获取分类的名称
func catKeyName(catKey int) string {
	switch catKey {
	case 1:
		return "软件类"
	case 2:
		return "硬件类"
	case 3:
		return "网络类"
	case 4:
		return "信息系统类"
	case 5:
		return "制造类"
	default:
		return ""
	}
}

// 辅助函数，用于获取子分类的名称
func subKeyName(catKey, subKey int) string {
	switch catKey {
	case 1:
		switch subKey {
		case 1:
			return "编程语言"
		case 2:
			return "前端开发"
		case 3:
			return "后端开发框架"
		case 4:
			return "数据库"
		case 5:
			return "容器技术"
		case 6:
			return "云计算"
		case 7:
			return "测试工具"
		case 8:
			return "版本控制"
		case 9:
			return "操作系统"
		case 10:
			return "移动开发"
		}
	case 2:
		switch subKey {
		case 1:
			return "嵌入式系统"
		case 2:
			return "硬件设计工具"
		case 3:
			return "微处理器"
		case 4:
			return "通信接口"
		case 5:
			return "传感器技术"
		case 6:
			return "电源管理"
		case 7:
			return "嵌入式软件开发工具"
		case 8:
			return "调试工具"
		case 9:
			return "生产工艺"
		case 10:
			return "驱动程序开发"
		}
	case 3:
		switch subKey {
		case 1:
			return "网络协议"
		case 2:
			return "网络设备"
		case 3:
			return "网络监控工具"
		case 4:
			return "网络虚拟化"
		case 5:
			return "网络安全"
		case 6:
			return "无线通信技术"
		case 7:
			return "网络操作系统"
		case 8:
			return "云网络"
	case 4:
		switch subKey {
		case 1:
			return "ERP系统"
		case 2:
			return "CRM系统"
		case 3:
			return "数据仓库"
		case 4:
			return "数据集成工具"
		case 5:
			return "商业智能"
		case 6:
			return "消息中间件"
		case 7:
			return "身份管理"
		case 8:
			return "大数据处理"
		case 9:
			return "企业应用集成"
		case 10:
			return "内容管理系统"
		}
	case 5:
		switch subKey {
		case 1:
			return "自动化控制"
		case 2:
			return "工业机器人"
		case 3:
			return "嵌入式系统"
		case 4:
			return "CAD/CAM软件"
		case 5:
			return "工业物联网"
		}
	}
	return ""
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

	// _, ok := data["subject_id"]
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
	// 	return
	// }

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

	data["subject_cat_key"], data["subject_sub_key"], datadata["subject_key"] = mapToID(data["subject_cat_key"], data["subject_sub_key"], datadata["subject_key"])

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
			log.Println("ChangePlan发生异常:", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	subject_id, ok := data["subject_id"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	_, ok = data["user_id"]
	if ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	study_time, ok := data["study_time"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	spend_time, ok := data["spend_time"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	addtime, ok := data["addtime"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	plan_id, ok := data["plan_id"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("ChangePlan:", "鉴权失败，用户不存在")
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}

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

	var planDetail Study

	db := baseClass.GetDB()
	db_result := db.Table("study_plans").Where("id = ?", plan_id).First(&planDetail)

	if planDetail.UserID != int(userID.(uint32)) {
		log.Println("ChangePlan:", "用户无权限访问该数据, planDetail:", planDetail.UserID, "userID:", userID)
		c.JSON(http.StatusOK, gin.H{"error": "用户无权限访问该数据"})
		return
	}

	// c.JSON(http.StatusOK, planDetail)
	// return

	if db_result.Error != nil {
		log.Println("ChangePlan数据库出错:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
		return
	}

	if val, ok := data["plan_name"].(string); ok {
		planDetail.PlanName = val
	}
	log.Println("1")
	planDetail.SubjectID = int(subject_id.(float64))
	log.Println("2")
	planDetail.StudyTime = int(study_time.(float64))

	planDetail.Spend_Time = spend_time.(float64)

	if val, ok := data["note"].(string); ok {
		planDetail.Note = val
	}
	if val, ok := data["tags"].(string); ok {
		planDetail.Tags = val
	}
	log.Println("3")
	planDetail.AddTime = int(addtime.(float64))

	db_result = db.Save(&planDetail)
	if db_result.Error != nil {
		log.Println("ChangePlan数据库出错:", db_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "success"})
	// if (data[])
}
