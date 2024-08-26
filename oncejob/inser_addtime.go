package oncejob

import (
	"CareerAnalysis/baseClass"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type Comment struct {
	ID      uint `gorm:"primaryKey"`
	AddTime int
}

func InserAddTime(c *gin.Context) {

	db := baseClass.GetDB()
	// 获取当前时间和两个月前的时间
	now := time.Now()
	twoMonthsAgo := now.AddDate(0, -2, 0)

	// 获取所有comments
	var comments []Comment
	if err := db.Find(&comments).Error; err != nil {
		panic("查询失败")
	}

	// 更新每个comment的addtime字段
	for _, comment := range comments {
		// 生成随机的时间戳
		randomTime := randate(twoMonthsAgo, now)

		randomTimestamp := int(randomTime.Unix())

		// 更新AddTime字段
		if err := db.Model(&comment).Update("addtime", randomTimestamp).Error; err != nil {
			fmt.Printf("更新记录 %d 失败: %v\n", comment.ID, err)
		} else {
			fmt.Printf("记录 %d 更新成功: %s\n", comment.ID, randomTime.Format("2006-01-02 15:04:05"))
		}
	}
}

// randate 生成在 start 和 end 之间的随机时间
func randate(start, end time.Time) time.Time {
	delta := end.Sub(start)
	sec := rand.Int63n(int64(delta.Seconds()))
	return start.Add(time.Duration(sec) * time.Second)
}
