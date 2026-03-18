package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 定义限流器接口。
//
// Allow 返回是否允许该 key 的一次请求。
type RateLimiter interface {
	Allow(key string) bool
}

// TokenBucketLimiter 是令牌桶限流器实现。
type TokenBucketLimiter struct {
	capacity int
	tokens   map[string]*bucket
	mu       sync.RWMutex
	rate     time.Duration
}

type bucket struct {
	tokens     int
	lastRefill time.Time
}

// NewTokenBucketLimiter 创建一个令牌桶限流器。
//
// capacity 为桶容量；rate 为令牌补充间隔（如 1s 表示每秒补 1 个令牌）。
func NewTokenBucketLimiter(capacity int, rate time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		capacity: capacity,
		tokens:   make(map[string]*bucket),
		rate:     rate,
	}
	go limiter.cleanup()
	return limiter
}

// Allow 判断是否允许该 key 的请求。
func (l *TokenBucketLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, exists := l.tokens[key]
	if !exists {
		b = &bucket{
			tokens:     l.capacity,
			lastRefill: time.Now(),
		}
		l.tokens[key] = b
	}

	now := time.Now()
	elapsed := now.Sub(b.lastRefill)
	tokensToAdd := int(elapsed / l.rate)

	if tokensToAdd > 0 {
		b.tokens = min(b.tokens+tokensToAdd, l.capacity)
		b.lastRefill = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

func (l *TokenBucketLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for key, b := range l.tokens {
			if now.Sub(b.lastRefill) > 30*time.Minute {
				delete(l.tokens, key)
			}
		}
		l.mu.Unlock()
	}
}

type RateLimitConfig struct {
	Limiter       RateLimiter
	KeyFunc       func(*gin.Context) string
	ExcludedPaths []string
	Message       string
}

// DefaultRateLimitConfig 返回默认限流配置。
//
// 默认优先按 user_id 限流，否则按 IP；默认排除 /health 和 /ping。
func DefaultRateLimitConfig(limiter RateLimiter) RateLimitConfig {
	return RateLimitConfig{
		Limiter: limiter,
		KeyFunc: func(c *gin.Context) string {
			if userID, exists := c.Get("user_id"); exists {
				return fmt.Sprintf("user:%v", userID)
			}
			return c.ClientIP()
		},
		ExcludedPaths: []string{"/health", "/ping"},
		Message:       "Too many requests",
	}
}

// RateLimit 创建通用限流中间件。
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range config.ExcludedPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		key := config.KeyFunc(c)
		if !config.Limiter.Allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": config.Message,
			})
			return
		}

		c.Next()
	}
}

// RateLimitByIP 创建按 IP 限流的中间件。
func RateLimitByIP(limiter RateLimiter, excludedPaths []string) gin.HandlerFunc {
	config := DefaultRateLimitConfig(limiter)
	config.KeyFunc = func(c *gin.Context) string {
		return c.ClientIP()
	}
	config.ExcludedPaths = excludedPaths
	return RateLimit(config)
}

// RateLimitByUserID 创建按 user_id 限流的中间件。
func RateLimitByUserID(limiter RateLimiter, excludedPaths []string) gin.HandlerFunc {
	config := DefaultRateLimitConfig(limiter)
	config.KeyFunc = func(c *gin.Context) string {
		if userID, exists := c.Get("user_id"); exists {
			return fmt.Sprintf("user:%v", userID)
		}
		return c.ClientIP()
	}
	config.ExcludedPaths = excludedPaths
	return RateLimit(config)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
