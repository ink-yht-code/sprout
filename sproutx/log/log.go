package log

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
	sugarLogger   *zap.SugaredLogger

	sinksMu sync.RWMutex
	sinks   []Sink

	defaultOnce sync.Once
)

// Sink 是可插拔日志接收器。
//
// 你可以实现该接口，把日志投递到 MQ/日志服务，再由下游入库。
// Sink 会接收到结构化日志 entry 和 fields。
type Sink interface {
	Write(entry zapcore.Entry, fields []zapcore.Field) error
	Sync() error
}

// RegisterSink 注册一个日志 Sink。
//
// 注意：当前实现只会在 Init 初始化时将 sinks 接入 logger。
func RegisterSink(s Sink) {
	if s == nil {
		return
	}
	sinksMu.Lock()
	defer sinksMu.Unlock()
	sinks = append(sinks, s)
}

// Config 是日志初始化配置。
type Config struct {
	Level    string `yaml:"level" json:"level" mapstructure:"level"`
	Encoding string `yaml:"encoding" json:"encoding" mapstructure:"encoding"`
	Output   string `yaml:"output" json:"output" mapstructure:"output"`

	Console *ConsoleConfig `yaml:"console" json:"console" mapstructure:"console"`
	File    *FileConfig    `yaml:"file" json:"file" mapstructure:"file"`
}

type ConsoleConfig struct {
	Enabled  *bool  `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	Encoding string `yaml:"encoding" json:"encoding" mapstructure:"encoding"`
	Output   string `yaml:"output" json:"output" mapstructure:"output"`
}

type FileConfig struct {
	Enabled  *bool  `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	Path     string `yaml:"path" json:"path" mapstructure:"path"`
	Encoding string `yaml:"encoding" json:"encoding" mapstructure:"encoding"`
}

// Init 初始化默认日志。
//
// Level 支持 debug/info/warn/error；Encoding 支持 json/console；Output 为输出位置（如 stdout 或文件路径）。
func Init(cfg Config) error {
	level := getLevel(cfg.Level)
	resolved := resolveOutputs(cfg)

	cores := make([]zapcore.Core, 0, 2)

	if resolved.consoleEnabled {
		encCfg := zap.NewProductionEncoderConfig()
		if resolved.consoleEncoding == "console" {
			encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		cores = append(cores, zapcore.NewCore(
			newEncoder(encCfg, resolved.consoleEncoding),
			zapcore.AddSync(stdSyncer(resolved.consoleOutput)),
			level,
		))
	}

	if resolved.fileEnabled {
		file, err := os.OpenFile(resolved.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}

		encCfg := zap.NewProductionEncoderConfig()
		cores = append(cores, zapcore.NewCore(
			newEncoder(encCfg, resolved.fileEncoding),
			zapcore.AddSync(file),
			level,
		))
	}

	if len(cores) == 0 {
		encCfg := zap.NewProductionEncoderConfig()
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encCfg),
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	logger := newLoggerWithSinks(zapcore.NewTee(cores...), level)
	defaultLogger = logger
	sugarLogger = logger.Sugar()
	return nil
}

type resolvedOutputs struct {
	consoleEnabled  bool
	consoleEncoding string
	consoleOutput   string

	fileEnabled  bool
	fileEncoding string
	filePath     string
}

func resolveOutputs(cfg Config) resolvedOutputs {
	hasNewCfg := cfg.Console != nil || cfg.File != nil

	if hasNewCfg {
		r := resolvedOutputs{}

		if cfg.Console != nil {
			r.consoleEnabled = cfg.Console.Enabled == nil || *cfg.Console.Enabled
			r.consoleEncoding = cfg.Console.Encoding
			r.consoleOutput = cfg.Console.Output
		}
		if cfg.File != nil {
			r.fileEnabled = cfg.File.Enabled == nil || *cfg.File.Enabled
			r.fileEncoding = cfg.File.Encoding
			r.filePath = cfg.File.Path
		}

		if r.consoleEncoding == "" {
			r.consoleEncoding = "console"
		}
		if r.consoleOutput == "" {
			r.consoleOutput = "stdout"
		}
		if r.fileEncoding == "" {
			r.fileEncoding = "json"
		}
		if r.filePath == "" {
			r.fileEnabled = false
		}

		return r
	}

	levelEncoding := cfg.Encoding
	if levelEncoding == "" {
		levelEncoding = "json"
	}
	output := cfg.Output
	if output == "" {
		output = "stdout"
	}

	if isStdoutLike(output) {
		return resolvedOutputs{
			consoleEnabled:  true,
			consoleEncoding: levelEncoding,
			consoleOutput:   output,
			fileEnabled:     false,
		}
	}

	return resolvedOutputs{
		consoleEnabled:  true,
		consoleEncoding: "console",
		consoleOutput:   "stdout",
		fileEnabled:     true,
		fileEncoding:    "json",
		filePath:        output,
	}
}

func newEncoder(encCfg zapcore.EncoderConfig, encoding string) zapcore.Encoder {
	switch encoding {
	case "console":
		return zapcore.NewConsoleEncoder(encCfg)
	default:
		return zapcore.NewJSONEncoder(encCfg)
	}
}

func stdSyncer(output string) *os.File {
	if strings.EqualFold(strings.TrimSpace(output), "stderr") {
		return os.Stderr
	}
	return os.Stdout
}

func newLoggerWithSinks(core zapcore.Core, level zapcore.Level) *zap.Logger {
	cores := []zapcore.Core{core}

	sinksMu.RLock()
	localSinks := append([]Sink(nil), sinks...)
	sinksMu.RUnlock()

	for _, s := range localSinks {
		cores = append(cores, newSinkCore(s, level))
	}

	return zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

type sinkCore struct {
	sink   Sink
	level  zapcore.LevelEnabler
	fields []zapcore.Field
}

func newSinkCore(s Sink, level zapcore.LevelEnabler) zapcore.Core {
	return &sinkCore{sink: s, level: level}
}

func (c *sinkCore) Enabled(lvl zapcore.Level) bool {
	return c.level.Enabled(lvl)
}

func (c *sinkCore) With(fields []zapcore.Field) zapcore.Core {
	if len(fields) == 0 {
		return c
	}
	cloned := *c
	cloned.fields = append(append([]zapcore.Field(nil), c.fields...), fields...)
	return &cloned
}

func (c *sinkCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *sinkCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	all := append(append([]zapcore.Field(nil), c.fields...), fields...)
	return c.sink.Write(ent, all)
}

func (c *sinkCore) Sync() error {
	return c.sink.Sync()
}

func isStdoutLike(output string) bool {
	switch strings.ToLower(strings.TrimSpace(output)) {
	case "stdout", "stderr":
		return true
	default:
		return false
	}
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
		defaultOnce.Do(func() {
			enabled := true
			_ = Init(Config{
				Level: "info",
				Console: &ConsoleConfig{
					Enabled:  &enabled,
					Encoding: "console",
					Output:   "stdout",
				},
			})
		})
	}
	return defaultLogger
}

// S 返回全局 zap.SugaredLogger（未初始化时返回 Nop logger）。
func S() *zap.SugaredLogger {
	if sugarLogger == nil {
		sugarLogger = L().Sugar()
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
