package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ink-yht-code/sprout/sproutx/log"
	"go.uber.org/zap"
)

// Config 是 HTTP Server 配置。
type Config struct {
	Enabled bool
	Addr    string
}

// Server 是 HTTP 服务封装，包含 gin.Engine 与 http.Server。
type Server struct {
	Engine *gin.Engine
	Server *http.Server
}

// NewServer 创建 HTTP Server。
func NewServer(cfg Config) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(Recovery(), RequestID(), Logger())

	return &Server{
		Engine: engine,
		Server: &http.Server{
			Addr:    cfg.Addr,
			Handler: engine,
		},
	}
}

// Run 启动 HTTP 服务（阻塞）。
func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}

// Shutdown 优雅关闭 HTTP 服务。
func (s *Server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// RequestID 为每个请求生成并注入 request_id。
//
// request_id 会写入 gin.Context 与响应头 X-Request-Id。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-Id", requestID)

		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// Logger 记录请求日志（method/path/status/latency 等）。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Ctx(c.Request.Context()).Info("HTTP request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		)
	}
}

// Recovery 捕获 panic 并返回 500。
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		log.Ctx(c.Request.Context()).Error("Panic recovered",
			zap.Any("error", err),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal error",
		})
	})
}
