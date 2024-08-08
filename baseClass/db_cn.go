package baseClass

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 结构体表示config.json中的配置
type Config struct {
	MySQL struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		DBName   string `json:"dbname"`
	} `json:"mysql"`
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	SSH struct {
		User       string `json:"user"`
		Host       string `json:"host"`
		Port       string `json:"port"`
		PrivateKey string `json:"private_key"`
	} `json:"ssh"`
	UseSSH bool `json:"use_ssh"` // 区分是否使用 SSH 隧道
}

// InitDB 初始化MySQL数据库连接
func InitDB() *gorm.DB {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil
	}

	dsn := config.MySQL.Username + ":" + config.MySQL.Password + "@tcp(" +
		config.MySQL.Host + ":" + config.MySQL.Port + ")/" + config.MySQL.DBName +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}

// InitRedis 初始化Redis连接
func InitRedis() *redis.Client {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password, // no password set
		DB:       config.Redis.DB,       // use default DB
	})
	return rdb
}

// CloseRedis 关闭Redis连接
func CloseRedis(rdb *redis.Client) {
	err := rdb.Close()
	if err != nil {
		log.Println("Failed to close Redis connection:", err)
	}
}

// GetRedisContext 获取Redis context
func GetRedisContext() context.Context {
	return context.Background()
}
