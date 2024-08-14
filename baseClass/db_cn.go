package baseClass

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"

	sql "github.com/go-sql-driver/mysql"
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

var db *gorm.DB

// InitDB 初始化MySQL数据库连接
func InitDB() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		os.Exit(1)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		os.Exit(1)
	}

	if !config.UseSSH {
		dsn := config.MySQL.Username + ":" + config.MySQL.Password + "@tcp(" +
			config.MySQL.Host + ":" + config.MySQL.Port + ")/" + config.MySQL.DBName +
			"?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
			os.Exit(1)
		}

	} else {
		// 使用SSH隧道连接数据库
		client := SSHConn()
		if client == nil {

			os.Exit(1)
		}

		// 定义符合 DialContextFunc 类型的拨号函数
		dial := func(ctx context.Context, addr string) (net.Conn, error) {
			// 修改 `addr` 参数，将其与 `network` 一起传递给 SSH 客户端的 Dial 方法
			// 这里我们将 `addr` 作为数据库主机和端口的地址
			return client.Dial("tcp", addr)
		}

		// 注册拨号函数到自定义网络 "mysql+ssh"
		sql.RegisterDialContext("mysql+ssh", sql.DialContextFunc(dial))

		// 使用自定义网络名 "mysql+ssh" 连接数据库
		dsn := config.MySQL.Username + ":" + config.MySQL.Password + "@mysql+ssh(" +
			config.MySQL.Host + ":" + config.MySQL.Port + ")/" + config.MySQL.DBName +
			"?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
			os.Exit(1)
		}

	}

}

func GetDB() *gorm.DB {
	return db
}

// // InitRedis 初始化Redis连接
// func InitRedis() *redis.Client {
// 	data, err := os.ReadFile("config.json")
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 		return nil
// 	}
// 	var config Config
// 	err = json.Unmarshal(data, &config)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 		return nil
// 	}
// 	if !config.UseSSH {
// 		// 不使用SSH直接连接Redis
// 		rdb := redis.NewClient(&redis.Options{
// 			Addr:     config.Redis.Addr,
// 			Password: config.Redis.Password, // no password set
// 			DB:       config.Redis.DB,       // use default DB
// 		})
// 		return rdb
// 	} else {
// 		// 使用SSH隧道连接Redis
// 		client := SSHConn()
// 		if client == nil {
// 			log.Println("Failed to establish SSH connection")
// 			return nil
// 		}

// 		// 定义一个自定义的拨号函数
// 		dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
// 			return client.Dial(network, addr)
// 		}

// 		// 创建 Redis 客户端，使用自定义的拨号器
// 		rdb := redis.NewClient(&redis.Options{
// 			Addr:     config.Redis.Addr,
// 			Password: config.Redis.Password, // no password set
// 			DB:       config.Redis.DB,       // use default DB
// 			Dialer:   dialer,                // 使用自定义拨号器
// 		})

// 		return rdb
// 	}
// }

// // CloseRedis 关闭Redis连接
// func CloseRedis(rdb *redis.Client) {
// 	err := rdb.Close()
// 	if err != nil {
// 		log.Println("Failed to close Redis connection:", err)
// 	}
// }

// // GetRedisContext 获取Redis context
// func GetRedisContext() context.Context {
// 	return context.Background()
// }
