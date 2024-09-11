package model

import (
	"net/http"

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
		for i, item := range request.Data {
			if i == 0 || i == 1 { // 前两个元素是数组
				// if array, ok := item.([]interface{}); ok {
				// 	for j, val := range array {
				// 		if
				// 	}
				// } else {
				// 	fmt.Printf("Element %d is expected to be an array, but is not.\n", i)
				// }
			} else { // 后续元素是单个值

			}
		}
	}
}
