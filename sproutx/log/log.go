package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
	sugarLogger   *zap.SugaredLogger
)

// Config 是日志初始化配置。
type Config struct {
	Level    string
	Encoding string
	Output   string
}

// Init 初始化默认日志。
//
// Level 支持 debug/info/warn/error；Encoding 支持 json/console；Output 为输出位置（如 stdout 或文件路径）。
func Init(cfg Config) error {
	level := getLevel(cfg.Level)
	encoding := cfg.Encoding
	if encoding == "" {
		encoding = "json"
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         encoding,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{cfg.Output},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	defaultLogger = logger
	sugarLogger = logger.Sugar()
	return nil
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// L 返回全局 zap.Logger（未初始化时返回 Nop logger）。
func L() *zap.Logger {
	if defaultLogger == nil {
		defaultLogger = zap.NewNop()
	}
	return defaultLogger
}

// S 返回全局 zap.SugaredLogger（未初始化时返回 Nop logger）。
func S() *zap.SugaredLogger {
	if sugarLogger == nil {
		sugarLogger = zap.NewNop().Sugar()
	}
	return sugarLogger
}

type ctxKey struct{}

// WithContext 将 logger 写入 context，便于链路传递。
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// Ctx 从 context 中取 logger。
//
// 若 context 中没有显式 logger，则尝试注入 request_id 字段。
func Ctx(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return L()
	}

	if logger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return logger
	}

	if requestID := GetRequestID(ctx); requestID != "" {
		return L().With(zap.String("request_id", requestID))
	}

	return L()
}

// GetRequestID 从 context 中读取 request_id。
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	for _, key := range []string{"request_id", "requestId", "x-request-id", "X-Request-Id"} {
		if v := ctx.Value(key); v != nil {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// Sync 刷新日志缓冲区。
func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Sync()
	}
	return nil
}

// Debug 输出 debug 级别日志。
func Debug(msg string, fields ...zap.Field) {
	L().Debug(msg, fields...)
}

// Info 输出 info 级别日志。
func Info(msg string, fields ...zap.Field) {
	L().Info(msg, fields...)
}

// Warn 输出 warn 级别日志。
func Warn(msg string, fields ...zap.Field) {
	L().Warn(msg, fields...)
}

// Error 输出 error 级别日志。
func Error(msg string, fields ...zap.Field) {
	L().Error(msg, fields...)
}

// Fatal 输出 fatal 级别日志并退出。
func Fatal(msg string, fields ...zap.Field) {
	L().Fatal(msg, fields...)
}

// DebugCtx 使用 context 中的 logger 输出 debug 级别日志。
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	Ctx(ctx).Debug(msg, fields...)
}

// InfoCtx 使用 context 中的 logger 输出 info 级别日志。
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	Ctx(ctx).Info(msg, fields...)
}

// WarnCtx 使用 context 中的 logger 输出 warn 级别日志。
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	Ctx(ctx).Warn(msg, fields...)
}

// ErrorCtx 使用 context 中的 logger 输出 error 级别日志。
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	Ctx(ctx).Error(msg, fields...)
}
