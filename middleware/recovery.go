package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/ink-yht-code/sprout/sproutx"
	"go.uber.org/zap"
)

// RecoveryConfig 是 Recovery（panic 恢复）中间件配置。
type RecoveryConfig struct {
	StackAll       bool
	StackSize      int
	DisablePrint   bool
	DisableStack   bool
	CustomRecovery func(c *gin.Context, err interface{})
}

// DefaultRecoveryConfig 返回默认 Recovery 配置。
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		StackAll:     false,
		StackSize:    4 << 10,
		DisablePrint: false,
		DisableStack: false,
	}
}

// Recovery 创建 Recovery 中间件。
//
// 捕获 panic 后记录日志与堆栈（可配置），并返回 500。
func Recovery(config RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if !config.DisablePrint {
					sproutx.L().Error("panic recovered",
						zap.Any("error", err),
						zap.String("path", c.Request.URL.Path),
						zap.String("method", c.Request.Method),
						zap.String("ip", c.ClientIP()),
					)
				}

				if !config.DisableStack {
					stack := debug.Stack()
					if config.StackAll {
						sproutx.L().Error("stack trace", zap.ByteString("stack", stack))
					} else {
						if len(stack) > config.StackSize {
							stack = stack[:config.StackSize]
						}
						sproutx.L().Error("stack trace", zap.ByteString("stack", stack))
					}
				}

				if config.CustomRecovery != nil {
					config.CustomRecovery(c, err)
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    http.StatusInternalServerError,
						"message": "Internal Server Error",
					})
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}

// RecoveryWithLogger 创建使用指定 zap.Logger 输出的 Recovery 中间件。
func RecoveryWithLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
					zap.ByteString("stack", debug.Stack()),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RecoveryWithWriter 创建将 panic 信息写入 writer 的 Recovery 中间件。
func RecoveryWithWriter(writer http.ResponseWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(writer, "panic: %v\n", err)
				fmt.Fprintf(writer, "%s\n", debug.Stack())

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
