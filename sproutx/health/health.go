package health

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HTTPHandler 返回一个简单的 HTTP 健康检查处理器。
func HTTPHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "sprout",
			"description": "基于 Gin 构建的 Go 微服务框架，提供路由、认证、校验等开箱即用的功能",
			"version":     "1.0.0",
		})
	}
}

// ReadyHandler 返回 readiness 处理器。
//
// checks 任意一个返回 false，则返回 503；全部通过则返回 200。
func ReadyHandler(checks ...func() bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, check := range checks {
			if !check() {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "not ready",
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	}
}

// Checker 定义依赖健康检查接口。
type Checker interface {
	Check(ctx context.Context) error
}

// HealthServer 是 gRPC 健康检查服务实现。
type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	checkers map[string]Checker
}

// NewHealthServer 创建一个 HealthServer。
func NewHealthServer() *HealthServer {
	return &HealthServer{checkers: make(map[string]Checker)}
}

// Register 注册某个 service 的健康检查器。
func (s *HealthServer) Register(service string, checker Checker) {
	s.checkers[service] = checker
}

// Check 实现 gRPC 健康检查接口。
//
// 当 req.Service 为空时返回整体 SERVING；当指定 service 未注册时返回 UNKNOWN。
func (s *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if req.Service == "" {
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
	}

	checker, ok := s.checkers[req.Service]
	if !ok {
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_UNKNOWN}, nil
	}

	if err := checker.Check(ctx); err != nil {
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

// Watch 实现 gRPC 健康检查的 Watch 接口。
//
// 当前实现为一次性返回 Check 结果。
func (s *HealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	resp, err := s.Check(stream.Context(), req)
	if err != nil {
		return err
	}
	return stream.Send(resp)
}

// DBChecker 是数据库健康检查器。
type DBChecker struct {
	ping func() error
}

// NewDBChecker 创建 DBChecker。
func NewDBChecker(ping func() error) *DBChecker {
	return &DBChecker{ping: ping}
}

// Check 执行数据库 ping。
func (c *DBChecker) Check(ctx context.Context) error {
	return c.ping()
}

// RedisChecker 是 Redis 健康检查器。
type RedisChecker struct {
	ping func() error
}

// NewRedisChecker 创建 RedisChecker。
func NewRedisChecker(ping func() error) *RedisChecker {
	return &RedisChecker{ping: ping}
}

// Check 执行 Redis ping。
func (c *RedisChecker) Check(ctx context.Context) error {
	return c.ping()
}
