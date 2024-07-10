package baseClass

import (
    "log"
    "github.com/go-redis/redis/v8"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "context"
)

// InitDB 初始化MySQL数据库连接
func InitDB() *gorm.DB {
    dsn := "root:123456@tcp(127.0.0.1:3306)/careeranalysis?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    return db
}

// InitRedis 初始化Redis连接
func InitRedis() *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
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
