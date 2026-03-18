package sprout

import (
	"net/http"
	"time"

	"github.com/ink-yht-code/sprout/context"
	"github.com/ink-yht-code/sprout/core"
	"github.com/ink-yht-code/sprout/jwt"
	"github.com/ink-yht-code/sprout/middleware"
	"github.com/ink-yht-code/sprout/session"
	"go.uber.org/zap"
)

// Context/Result 等类型别名用于对外暴露统一入口，隐藏内部包结构。
type (
	// Context 是增强版上下文，基于 gin.Context 封装。
	Context = context.Context
	// Result 是统一 JSON 响应结构。
	Result = core.Result

	// PageData 是分页返回结构。
	PageData[T any] core.PageData[T]
	// PageRequest 是分页请求参数。
	PageRequest = core.PageRequest

	// Session 表示当前请求的会话。
	Session = session.Session
	// Provider 是 Session 存储提供者。
	Provider = session.Provider
	// TokenCarrier 描述 Token 的携带方式（Cookie/Header）。
	TokenCarrier = session.TokenCarrier

	// Claims 是 JWT Claims。
	Claims = jwt.Claims
	// TokenPair 是 Access/Refresh token 对。
	TokenPair = jwt.TokenPair
	// Options 是 JWT 配置。
	Options = jwt.Options
	// Manager 是 JWT 管理器接口。
	Manager = jwt.Manager

	// CORSConfig 是 CORS 中间件配置。
	CORSConfig = middleware.CORSConfig
	// CSRFConfig 是 CSRF 中间件配置。
	CSRFConfig = middleware.CSRFConfig
	// RateLimitConfig 是限流中间件配置。
	RateLimitConfig = middleware.RateLimitConfig
	// RecoveryConfig 是 Recovery 中间件配置。
	RecoveryConfig = middleware.RecoveryConfig
	// JWTAuthConfig 是 JWT 鉴权中间件配置。
	JWTAuthConfig = middleware.JWTAuthConfig
)

const (
	// CodeSuccess 表示成功。
	CodeSuccess = core.CodeSuccess
	// CodeWarning 表示警告。
	CodeWarning = core.CodeWarning
	// CodeError 表示通用错误。
	CodeError = core.CodeError
	// CodeInvalidParam 表示参数错误。
	CodeInvalidParam = core.CodeInvalidParam
	// CodeInternalError 表示系统内部错误。
	CodeInternalError = core.CodeInternalError
	// CodeUnauthorized 表示未授权。
	CodeUnauthorized = core.CodeUnauthorized
	// CodeForbidden 表示无权限。
	CodeForbidden = core.CodeForbidden
	// CodeNotFound 表示资源不存在。
	CodeNotFound = core.CodeNotFound
	// CodeConflict 表示资源冲突。
	CodeConflict = core.CodeConflict
	// CodeTooManyRequests 表示请求过于频繁。
	CodeTooManyRequests = core.CodeTooManyRequests
	// CodeServiceUnavailable 表示服务不可用。
	CodeServiceUnavailable = core.CodeServiceUnavailable
)

var (
	// ErrNoResponse 表示不需要写回响应。
	ErrNoResponse = core.ErrNoResponse
	// ErrUnauthorized 表示未授权。
	ErrUnauthorized = core.ErrUnauthorized
	// ErrSessionNotFound 表示会话不存在。
	ErrSessionNotFound = core.ErrSessionNotFound
	// ErrSessionExpired 表示会话已过期。
	ErrSessionExpired = core.ErrSessionExpired
	// ErrInvalidToken 表示无效 Token。
	ErrInvalidToken = core.ErrInvalidToken
	// ErrProviderNotInitialized 表示 Session Provider 未初始化。
	ErrProviderNotInitialized = session.ErrProviderNotInitialized

	// CodeMessage 是默认错误码消息映射。
	CodeMessage = core.CodeMessage
)

// Context 是对 context.Context 的类型别名，便于 core 包复用增强 Context。
type Context = context.Context

// GetCodeMessage 获取业务码的默认中文提示。
func GetCodeMessage(code int) string {
	return core.GetCodeMessage(code)
}

// Success 构造成功响应。
func Success(msg string, data any) Result {
	return core.Success(msg, data)
}

// Warning 构造警告响应。
func Warning(msg string, data any) Result {
	return core.Warning(msg, data)
}

// Error 构造通用错误响应。
func Error(msg string) Result {
	return core.Error(msg)
}

// ErrorWithCode 构造指定业务码的错误响应。
func ErrorWithCode(code int, msg string) Result {
	return core.ErrorWithCode(code, msg)
}

// InvalidParam 构造参数错误响应。
func InvalidParam(msg string) Result {
	return core.InvalidParam(msg)
}

// InternalError 构造系统内部错误响应。
func InternalError() Result {
	return core.InternalError()
}

// Unauthorized 构造未授权响应。
func Unauthorized() Result {
	return core.Unauthorized()
}

// Forbidden 构造无权限响应。
func Forbidden() Result {
	return core.Forbidden()
}

// NotFound 构造资源不存在响应。
func NotFound(msg string) Result {
	return core.NotFound(msg)
}

// Conflict 构造资源冲突响应。
func Conflict(msg string) Result {
	return core.Conflict(msg)
}

// TooManyRequests 构造请求过于频繁响应。
func TooManyRequests() Result {
	return core.TooManyRequests()
}

// ServiceUnavailable 构造服务不可用响应。
func ServiceUnavailable() Result {
	return core.ServiceUnavailable()
}

// W 将业务函数包装为 Handler（实际类型为 gin.HandlerFunc）。
func W(fn func(ctx *context.Context) (Result, error)) interface{} {
	return core.W(fn)
}

// B 将业务函数包装为 Handler，并自动绑定请求参数。
func B[Req any](fn func(ctx *context.Context, req Req) (Result, error)) interface{} {
	return core.B(fn)
}

// S 将业务函数包装为 Handler，并自动获取 Session。
func S(fn func(ctx *context.Context, sess session.Session) (Result, error)) interface{} {
	return core.S(fn)
}

// BS 将业务函数包装为 Handler，同时支持参数绑定和 Session 获取。
func BS[Req any](fn func(ctx *context.Context, req Req, sess session.Session) (Result, error)) interface{} {
	return core.BS(fn)
}

// NewJWTManager 创建 JWT Manager。
func NewJWTManager(opts jwt.Options) jwt.Manager {
	return jwt.NewManager(opts)
}

// NewJWTOptions 创建 JWT Options（以秒为单位指定过期时间）。
func NewJWTOptions(signKey string, accessExpire, refreshExpire int64) jwt.Options {
	return jwt.NewOptions(signKey,
		time.Duration(accessExpire)*time.Second,
		time.Duration(refreshExpire)*time.Second)
}

// SetDefaultSessionProvider 设置默认 Session Provider。
func SetDefaultSessionProvider(provider session.Provider) {
	session.SetDefaultProvider(provider)
}

// GetSession 从 Context 中获取当前 Session。
func GetSession(ctx *context.Context) (session.Session, error) {
	return session.Get(ctx)
}

// NewSession 创建并写入新的 Session。
func NewSession(ctx *context.Context, userId string, jwtData map[string]string, sessData map[string]any) (session.Session, error) {
	return session.NewSession(ctx, userId, jwtData, sessData)
}

// DefaultCORSConfig 返回默认 CORS 配置。
func DefaultCORSConfig() middleware.CORSConfig {
	return middleware.DefaultCORSConfig()
}

// CORS 创建 CORS 中间件（实际类型为 gin.HandlerFunc）。
func CORS(config middleware.CORSConfig) interface{} {
	return middleware.CORS(config)
}

// DefaultCSRFConfig 返回默认 CSRF 配置。
func DefaultCSRFConfig() middleware.CSRFConfig {
	return middleware.DefaultCSRFConfig()
}

// CSRF 创建 CSRF 中间件（实际类型为 gin.HandlerFunc）。
func CSRF(config middleware.CSRFConfig) interface{} {
	return middleware.CSRF(config)
}

// NewTokenBucketLimiter 创建令牌桶限流器。
func NewTokenBucketLimiter(capacity int, rate time.Duration) middleware.RateLimiter {
	return middleware.NewTokenBucketLimiter(capacity, rate)
}

// DefaultRateLimitConfig 返回默认限流配置。
func DefaultRateLimitConfig(limiter middleware.RateLimiter) middleware.RateLimitConfig {
	return middleware.DefaultRateLimitConfig(limiter)
}

// RateLimit 创建限流中间件（实际类型为 gin.HandlerFunc）。
func RateLimit(config middleware.RateLimitConfig) interface{} {
	return middleware.RateLimit(config)
}

// RateLimitByIP 按客户端 IP 维度限流。
func RateLimitByIP(limiter middleware.RateLimiter, excludedPaths []string) interface{} {
	return middleware.RateLimitByIP(limiter, excludedPaths)
}

// RateLimitByUserID 按 user_id 维度限流。
func RateLimitByUserID(limiter middleware.RateLimiter, excludedPaths []string) interface{} {
	return middleware.RateLimitByUserID(limiter, excludedPaths)
}

// DefaultRecoveryConfig 返回默认 Recovery 配置。
func DefaultRecoveryConfig() middleware.RecoveryConfig {
	return middleware.DefaultRecoveryConfig()
}

// Recovery 创建 Recovery 中间件（实际类型为 gin.HandlerFunc）。
func Recovery(config middleware.RecoveryConfig) interface{} {
	return middleware.Recovery(config)
}

// RecoveryWithLogger 创建带 zap Logger 的 Recovery 中间件。
func RecoveryWithLogger(logger interface{}) interface{} {
	switch l := logger.(type) {
	case *zap.Logger:
		return middleware.RecoveryWithLogger(l)
	default:
		return middleware.Recovery(middleware.DefaultRecoveryConfig())
	}
}

// RecoveryWithWriter 创建带自定义 writer 的 Recovery 中间件。
func RecoveryWithWriter(writer interface{}) interface{} {
	if w, ok := writer.(http.ResponseWriter); ok {
		return middleware.RecoveryWithWriter(w)
	}
	return middleware.Recovery(middleware.DefaultRecoveryConfig())
}

// DefaultJWTAuthConfig 返回默认 JWT 鉴权配置。
func DefaultJWTAuthConfig(manager jwt.Manager) middleware.JWTAuthConfig {
	return middleware.DefaultJWTAuthConfig(manager)
}

// JWTAuth 创建 JWT 鉴权中间件（实际类型为 gin.HandlerFunc）。
func JWTAuth(config middleware.JWTAuthConfig) interface{} {
	return middleware.JWTAuth(config)
}
