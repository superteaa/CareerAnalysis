package model

import (
	"CareerAnalysis/baseClass"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// func SelectRole(c *gin.Context) {
// 	role := c.Query("role")
// 	if role == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 	}
// 	if role == "学生" {
// 		c.JSON(http.StatusOK, "ok")
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"error": "更多功能正在开发中"})
// 	}
// }

func DealQ(c *gin.Context) {
	var request struct {
		Is_Test int           `json:"is_test" binding:"required"`
		Data    []interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid request"})
		return
	}

	if request.Is_Test != 1 {
		c.JSON(http.StatusOK, "ok")
		return
	} else {
		// 处理 Data 字段
		jobs_map := make(map[int]bool) // 用于记录已添加的元素
		jobs_arr := []int{}            // 用于存储唯一的元素
		for i, item := range request.Data {
			if i == 0 || i == 1 {
				continue
			} else if i == 2 { // 后续元素是单个值
				switch item {
				case "A":
					if !jobs_map[3] {
						jobs_arr = append(jobs_arr, 3)
						jobs_map[3] = true
					}
				case "B":
					if !jobs_map[4] {
						jobs_arr = append(jobs_arr, 4)
						jobs_map[4] = true
					}
				case "C":
					if !jobs_map[1] {
						jobs_arr = append(jobs_arr, 1)
						jobs_map[1] = true
					}
				case "D":
					break
				default:
					c.JSON(http.StatusOK, gin.H{"error": "Invalid request"})
					return
				}
			} else if i == 3 {
				switch item {
				case "A":
					break
				case "B":
					if !jobs_map[1] {
						jobs_arr = append(jobs_arr, 1)
						jobs_map[1] = true
					}
					if !jobs_map[3] {
						jobs_arr = append(jobs_arr, 3)
						jobs_map[3] = true
					}
					if !jobs_map[4] {
						jobs_arr = append(jobs_arr, 4)
						jobs_map[4] = true
					}

				case "C":
					if !jobs_map[2] {
						jobs_arr = append(jobs_arr, 2)
						jobs_map[2] = true
					}

				case "D":
					if !jobs_map[1] {
						jobs_arr = append(jobs_arr, 1)
						jobs_map[1] = true
					}
				default:
					c.JSON(http.StatusOK, gin.H{"error": "Invalid request"})
					return
				}
			} else if i == 4 {
				switch item {
				case "A":
					if !jobs_map[3] {
						jobs_arr = append(jobs_arr, 3)
						jobs_map[3] = true
					}
				case "B":
					if !jobs_map[2] {
						jobs_arr = append(jobs_arr, 2)
						jobs_map[2] = true
					}
					if !jobs_map[4] {
						jobs_arr = append(jobs_arr, 4)
						jobs_map[4] = true
					}

				case "C":
					break
				case "D":
					break
				default:
					c.JSON(http.StatusOK, gin.H{"error": "Invalid request"})
					return
				}
			} else if i == 5 {
				switch item {
				case "A":
					if !jobs_map[1] {
						jobs_arr = append(jobs_arr, 1)
						jobs_map[1] = true
					}
					if !jobs_map[3] {
						jobs_arr = append(jobs_arr, 3)
						jobs_map[3] = true
					}
				case "B":

					if !jobs_map[4] {
						jobs_arr = append(jobs_arr, 4)
						jobs_map[4] = true
					}

				case "C":
					if !jobs_map[2] {
						jobs_arr = append(jobs_arr, 2)
						jobs_map[2] = true
					}
				case "D":
					break
				default:
					c.JSON(http.StatusOK, gin.H{"error": "Invalid request"})
					return
				}
			}
		}
		// 处理 tags
		// if tagsArr, ok := data["tags"].([]interface{}); ok {
		var strSlice []string
		for _, item := range jobs_arr {
			strSlice = append(strSlice, strconv.FormatInt(int64(item), 10))

		}
		jobsStr := strings.Join(strSlice, ",")

		userID, exists := c.Get("userID")
		if !exists {
			log.Println("DealQ:", "鉴权失败，用户不存在")
			c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
			return
		}

		recomment := map[string]interface{}{
			"job_arr": jobsStr,
			"user_id": int(userID.(uint32)),
		}

		// 启动事务
		db := baseClass.GetDB()
		tx := db.Begin()

		var existingRecord struct {
			Job_arr string
			User_id int
		}

		// 查询是否已存在相同的 user_id
		if err := tx.Table("recomment_jobs").Where("user_id = ?", recomment["user_id"]).First(&existingRecord).Error; err == nil {
			// 如果找到相同的 user_id
			log.Println("DealQ数据库出错: 已存在user_id", recomment["user_id"])
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "请勿重复提交"})
			return
		}

		// 插入数据
		if err := tx.Table("recomment_jobs").Create(&recomment).Error; err != nil {
			log.Println("DealQ数据库出错:", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统优化中"})
			return
		}

		// 提交事务
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
	}
}
