package middleware

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// ActiveLimitBuilder 用于构建活跃连接（并发请求）限制中间件。
type ActiveLimitBuilder struct {
	maxActive int64
}

// NewActiveLimitBuilder 创建一个 ActiveLimitBuilder。
//
// maxActive 表示最大并发请求数。
func NewActiveLimitBuilder(maxActive int64) *ActiveLimitBuilder {
	return &ActiveLimitBuilder{maxActive: maxActive}
}

// Build 构建活跃连接限制中间件。
//
// 当并发请求数超过 maxActive 时，直接返回 429。
func (b *ActiveLimitBuilder) Build() gin.HandlerFunc {
	var currentActive int64

	return func(c *gin.Context) {
		current := atomic.AddInt64(&currentActive, 1)
		defer func() {
			atomic.AddInt64(&currentActive, -1)
		}()

		if current > b.maxActive {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
