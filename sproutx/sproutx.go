package sproutx

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ink-yht-code/sprout/sproutx/app"
	"github.com/ink-yht-code/sprout/sproutx/db"
	xerror "github.com/ink-yht-code/sprout/sproutx/error"
	"github.com/ink-yht-code/sprout/sproutx/health"
	"github.com/ink-yht-code/sprout/sproutx/httpx"
	"github.com/ink-yht-code/sprout/sproutx/log"
	"github.com/ink-yht-code/sprout/sproutx/outbox"
	"github.com/ink-yht-code/sprout/sproutx/redis"
	"github.com/ink-yht-code/sprout/sproutx/rpc"
	"github.com/ink-yht-code/sprout/sproutx/tx"
	redislib "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App/Config 等类型别名用于对外暴露统一入口，隐藏内部包结构。
type (
	// App 是应用启动器。
	App = app.App
	// Config 是应用配置。
	Config = app.Config

	// ServiceConfig 是服务基础配置。
	ServiceConfig = app.ServiceConfig

	// HTTPConfig 是 HTTP 服务配置。
	HTTPConfig = httpx.Config
	// Server 是 HTTP Server。
	Server = httpx.Server

	// LogConfig 是日志配置。
	LogConfig = log.Config

	// DBConfig 是数据库配置。
	DBConfig = db.Config

	// RedisConfig 是 Redis 配置。
	RedisConfig = redis.Config

	// TxManager 是事务管理器。
	TxManager = tx.Manager

	// GRPCConfig 是 gRPC 服务配置。
	GRPCConfig = rpc.Config
	// GRPCServer 是 gRPC Server。
	GRPCServer = rpc.Server
	// GRPCClient 是 gRPC Client。
	GRPCClient = rpc.Client

	// HealthServer 是健康检查 HTTP Server。
	HealthServer = health.HealthServer
	// HealthChecker 是健康检查器。
	HealthChecker = health.Checker
	// DBChecker 是数据库健康检查。
	DBChecker = health.DBChecker
	// RedisChecker 是 Redis 健康检查。
	RedisChecker = health.RedisChecker

	// Outbox 是事务消息 Outbox。
	Outbox = outbox.Outbox
	// OutboxRepo 是 Outbox 仓储（轻量封装）。
	OutboxRepo = outbox.Repo
	// OutboxRepository 是 Outbox Repository 接口。
	OutboxRepository = outbox.Repository

	// BizError 是业务错误类型。
	BizError = xerror.BizError
)

// NewApp 创建应用启动器。
func NewApp(cfg *Config) (*App, error) {
	return app.New(cfg)
}

// NewHTTPServer 创建 HTTP Server。
func NewHTTPServer(cfg HTTPConfig) *Server {
	return httpx.NewServer(cfg)
}

// NewDB 创建 GORM DB。
func NewDB(cfg DBConfig) (*gorm.DB, error) {
	return db.New(cfg)
}

// NewRedis 创建 Redis 客户端。
func NewRedis(cfg RedisConfig) *redislib.Client {
	return redis.New(cfg)
}

// NewTxManager 创建事务管理器。
func NewTxManager(db *gorm.DB) *TxManager {
	return tx.NewManager(db)
}

// NewGRPCServer 创建 gRPC Server。
func NewGRPCServer(cfg GRPCConfig) *GRPCServer {
	return rpc.NewServer(cfg)
}

// NewGRPCClient 创建 gRPC Client。
func NewGRPCClient(addr string) (*GRPCClient, error) {
	return rpc.NewClient(addr)
}

// NewHealthServer 创建健康检查 Server。
func NewHealthServer() *HealthServer {
	return health.NewHealthServer()
}

// NewDBChecker 创建数据库健康检查器。
func NewDBChecker(ping func() error) *DBChecker {
	return health.NewDBChecker(ping)
}

// NewRedisChecker 创建 Redis 健康检查器。
func NewRedisChecker(ping func() error) *RedisChecker {
	return health.NewRedisChecker(ping)
}

// NewOutboxRepo 创建 Outbox Repo。
func NewOutboxRepo() *OutboxRepo {
	return outbox.NewRepo()
}

// MapBizErrorToHTTP 将业务错误映射为 HTTP 响应。
func MapBizErrorToHTTP(c *gin.Context, err error) {
	xerror.MapToHTTP(c, err)
}

// ErrorHandler 返回业务错误处理中间件。
func ErrorHandler() gin.HandlerFunc {
	return xerror.Handler()
}

// InitLog 初始化 zap 日志。
func InitLog(cfg LogConfig) error {
	return log.Init(cfg)
}

// L 返回全局 zap.Logger。
func L() *zap.Logger {
	return log.L()
}

// S 返回全局 zap.SugaredLogger。
func S() *zap.SugaredLogger {
	return log.S()
}

// Ctx 从 context 获取带 request_id 的 logger。
func Ctx(ctx context.Context) *zap.Logger {
	return log.Ctx(ctx)
}

// GetDBFromCtx 从 context 获取当前事务内的 DB（若无事务则返回 defaultDB）。
func GetDBFromCtx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	return tx.FromContext(ctx, defaultDB)
}

// GetRequestID 从 context 中提取 request_id。
func GetRequestID(ctx context.Context) string {
	return log.GetRequestID(ctx)
}

// SyncLog 刷新并落盘日志。
func SyncLog() error {
	return log.Sync()
}
