package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Config 是中间件整体配置集合。
//
// 你可以按需启用/配置 CORS、CSRF、限流、Recovery 等中间件。
type Config struct {
	CORS        CORSConfig
	CSRF        CSRFConfig
	RateLimit   RateLimitConfig
	Recovery    RecoveryConfig
	AccessLog   any
	ActiveLimit any
}

// DefaultConfig 返回一套默认中间件配置。
func DefaultConfig() Config {
	return Config{
		CORS:        DefaultCORSConfig(),
		CSRF:        DefaultCSRFConfig(),
		RateLimit:   DefaultRateLimitConfig(nil),
		Recovery:    DefaultRecoveryConfig(),
		AccessLog:   nil,
		ActiveLimit: nil,
	}
}

// NewDefaultMiddleware 返回默认建议启用的中间件列表。
//
// 默认包含 Recovery 与 CORS。
func NewDefaultMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		Recovery(DefaultRecoveryConfig()),
		CORS(DefaultCORSConfig()),
	}
}

// NewRateLimitMiddleware 创建按 IP 限流的中间件。
//
// requestsPerMinute 为每分钟允许的请求数，默认排除 /health 和 /ping。
func NewRateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewTokenBucketLimiter(requestsPerMinute, time.Minute)
	return RateLimitByIP(limiter, []string{"/health", "/ping"})
}
