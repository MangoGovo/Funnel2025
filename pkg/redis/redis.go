package redis

import (
	"funnel/pkg/config"
	_ "funnel/pkg/log"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// redisConfig 定义 Redis 数据库的配置结构体
type redisConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
}

// GlobalClient 全局 Redis 客户端实例
var GlobalClient *redis.Client

// init 函数用于初始化 Redis 客户端和配置信息
func init() {
	info := redisConfig{
		Host:     config.Config.GetString("redis.host"),
		Port:     config.Config.GetString("redis.port"),
		DB:       config.Config.GetInt("redis.db"),
		Password: config.Config.GetString("redis.pass"),
	}

	GlobalClient = redis.NewClient(&redis.Options{
		Addr:     info.Host + ":" + info.Port,
		Password: info.Password,
		DB:       info.DB,
	})

	zap.L().Info("Redis初始化成功")
}
