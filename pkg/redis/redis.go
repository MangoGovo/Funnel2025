package redis

import (
	"context"

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

var (
	// Client 全局 Redis 客户端实例
	Client *redis.Client
	// Ctx 全局 redis 上下文
	Ctx = context.Background()
)

// init 函数用于初始化 Redis 客户端和配置信息
func init() {
	info := redisConfig{
		Host:     config.Config.GetString("redis.host"),
		Port:     config.Config.GetString("redis.port"),
		DB:       config.Config.GetInt("redis.db"),
		Password: config.Config.GetString("redis.pass"),
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     info.Host + ":" + info.Port,
		Password: info.Password,
		DB:       info.DB,
	})
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		zap.L().Info("Redis初始化失败", zap.Error(err))
		return
	}
	zap.L().Info("Redis初始化成功")
}
