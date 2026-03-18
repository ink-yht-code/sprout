package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// Config 是 Redis 客户端配置。
type Config struct {
	Addr     string
	Password string
	DB       int
}

// New 创建 Redis 客户端。
func New(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}
