package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewNewCommand 创建 new 子命令。
//
// 用于创建新的服务项目骨架（包含目录结构、配置、示例 API 与 wiring）。
func NewNewCommand() *cobra.Command {
	var (
		transportFlag string
		daoFlag       string
		cacheFlag     string
		registryFlag  string
	)

	cmd := &cobra.Command{
		Use:   "new service <name>",
		Short: "创建新服务",
		Long:  `创建一个新的微服务项目，包含完整的目录结构和基础文件`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[1]
			if err := createService(serviceName, transportFlag, daoFlag, cacheFlag, registryFlag); err != nil {
				color.Red("创建服务失败: %v", err)
				os.Exit(1)
			}
			color.Green("✓ 服务 %s 创建成功！", serviceName)
			color.Yellow("\n下一步:")
			fmt.Printf("  cd %s\n", serviceName)
			fmt.Println("  编辑 .sprout 文件定义 API")
			fmt.Println("  运行: sprout-gen api " + serviceName)
		},
	}

	cmd.Flags().StringVarP(&transportFlag, "transport", "t", "http", "传输协议: http, rpc, http,rpc")
	cmd.Flags().StringVar(&daoFlag, "dao", "gorm", "DAO 类型")
	cmd.Flags().StringVar(&cacheFlag, "cache", "redis", "Cache 类型")
	cmd.Flags().StringVarP(&registryFlag, "registry", "r", "http://localhost:18080", "Registry 服务地址")
	return cmd
}

func createService(serviceName string, transport string, dao string, cache string, registryAddr string) error {
	serviceDir := serviceName
	if _, err := os.Stat(serviceDir); !os.IsNotExist(err) {
		return fmt.Errorf("目录 %s 已存在", serviceDir)
	}

	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	dirs := []string{
		"configs",
		"cmd",
		"internal/config",
		"internal/domain/entity",
		"internal/domain/port",
		"internal/domain/errs",
		"internal/domain/event",
		"internal/repository/dao",
		"internal/service",
		"internal/types",
		"internal/web",
		"internal/wiring",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(serviceDir, dir), 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}

	serviceID := 101
	if registryAddr != "" {
		id, err := allocateServiceID(registryAddr, serviceName)
		if err == nil {
			serviceID = id
		}
	}

	files := map[string]string{
		"go.mod":                             goModTemplate,
		".sprout":                            sproutTemplate,
		"README.md":                          readmeTemplate,
		"configs/config.yaml":                configTemplate,
		"cmd/main.go":                        mainTemplate,
		"internal/config/config.go":          configStructTemplate,
		"internal/domain/entity/entity.go":   entityTemplate,
		"internal/domain/port/repository.go": repositoryPortTemplate,
		"internal/domain/errs/codes.go":      codesTemplate,
		"internal/domain/errs/error.go":      errorTemplate,
		"internal/service/service.go":        serviceTemplate,
		"internal/types/types.go":            typesTemplate,
		"internal/web/handler.go":            handlerTemplate,
		"internal/wiring/wiring.go":          wiringTemplate,
	}

	for path, tmpl := range files {
		if err := createFileFromTemplate(serviceDir, path, tmpl, serviceName, serviceID, transport); err != nil {
			return fmt.Errorf("创建文件 %s 失败: %w", path, err)
		}
	}

	_ = dao
	_ = cache
	return nil
}

func createFileFromTemplate(baseDir, path, tmpl string, serviceName string, serviceID int, transport string) error {
	fullPath := filepath.Join(baseDir, path)

	t, err := template.New(path).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	data := struct {
		ServiceName string
		ServiceID   int
		Transport   string
	}{
		ServiceName: serviceName,
		ServiceID:   serviceID,
		Transport:   transport,
	}

	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

func allocateServiceID(registryAddr string, serviceName string) (int, error) {
	addr := strings.TrimSpace(registryAddr)
	if addr == "" {
		return 0, fmt.Errorf("registry address is empty")
	}
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}
	addr = strings.TrimRight(addr, "/")

	reqBody, err := json.Marshal(map[string]string{"name": serviceName})
	if err != nil {
		return 0, fmt.Errorf("marshal allocate request: %w", err)
	}

	url := addr + "/v1/services:allocate"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request registry allocate: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("registry allocate failed: status=%d body=%s", resp.StatusCode, string(respBytes))
	}

	var out struct {
		ServiceID int    `json:"service_id"`
		Name      string `json:"name"`
	}
	if err := json.Unmarshal(respBytes, &out); err != nil {
		return 0, fmt.Errorf("decode registry response: %w; body=%s", err, string(respBytes))
	}
	if out.ServiceID <= 0 {
		return 0, fmt.Errorf("invalid service_id from registry: %d", out.ServiceID)
	}
	return out.ServiceID, nil
}

const goModTemplate = `module {{.ServiceName}}
 
 go 1.25.0
 
 require (
	github.com/ink-yht-code/sprout v0.0.2
 )
`

const sproutTemplate = `syntax = "v1"

type ExampleReq {
    Name string ` + "`" + `json:"name" validate:"required"` + "`" + `
}

type ExampleResp {
    Message string ` + "`" + `json:"message"` + "`" + `
}

server {
    prefix "/api"
}

service {{.ServiceName}} {
    public {
        GET "/ping" Ping -> ExampleResp
    }
    
    private {
    }
}
`

const readmeTemplate = `# {{.ServiceName}} Service

{{.ServiceName}} 微服务，基于 Sprout 框架构建。

## 快速开始

安装依赖:
go mod download

运行服务:
go run cmd/main.go

## API 文档

### 公开接口

- GET /api/ping - 健康检查

### 私有接口

需要认证的接口

## 配置

配置文件位于 configs/config.yaml。

## 开发

生成代码:
sprout-gen api {{.ServiceName}}

运行测试:
go test ./...
`

const configTemplate = `service:
  id: {{.ServiceID}}
  name: {{.ServiceName}}

http:
  enabled: true
  addr: ":8080"

grpc:
  enabled: false
  addr: ":9090"

log:
  # 日志级别: debug|info|warn|error
  level: "info"
  # 旧版兼容配置: json|console
  # - output 为 stdout/stderr 时：按 encoding 输出到终端
  # - output 为文件路径时：默认启用 "终端 console + 文件 json" 的双输出
  encoding: "json"
  # 旧版兼容配置: stdout|stderr|文件路径
  output: "stdout"

  # 推荐：分别配置终端与文件输出（更直观）
  # - 终端输出用 console（可读）
  # - 文件输出用 json（便于采集入库）
  console:
    # 是否启用终端输出
    enabled: true
    # console|json
    encoding: "console"
    # stdout|stderr
    output: "stdout"
  file:
    # 是否启用文件输出
    enabled: false
    # 日志文件路径
    path: "logs/app.log"
    # 建议保持 json，便于结构化采集
    encoding: "json"

db:
  dsn: ""
  max_open: 100
  max_idle: 10
  log_level: "info"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
`

const mainTemplate = `package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ink-yht-code/{{.ServiceName}}/internal/config"
	"github.com/ink-yht-code/{{.ServiceName}}/internal/wiring"
	"github.com/ink-yht-code/sproutx"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	application, err := wiring.InitApp(cfg)
	if err != nil {
		panic(fmt.Sprintf("初始化应用失败: %v", err))
	}

	go func() {
		sproutx.L().Info("服务启动中", 
			zap.String("service", cfg.Service.Name),
			zap.Int("service_id", cfg.Service.ID))
		if err := application.Run(); err != nil {
			sproutx.L().Fatal("服务运行失败", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sproutx.L().Info("正在关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		sproutx.L().Error("关闭服务失败", zap.Error(err))
	}

	sproutx.L().Info("服务已关闭")
}
`

const configStructTemplate = `package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service ServiceConfig ` + "`" + `yaml:"service"` + "`" + `
	HTTP    HTTPConfig    ` + "`" + `yaml:"http"` + "`" + `
	GRPC    GRPCConfig    ` + "`" + `yaml:"grpc"` + "`" + `
	Log     LogConfig     ` + "`" + `yaml:"log"` + "`" + `
	DB      DBConfig      ` + "`" + `yaml:"db"` + "`" + `
	Redis   RedisConfig   ` + "`" + `yaml:"redis"` + "`" + `
}

type ServiceConfig struct {
	ID   int    ` + "`" + `yaml:"id"` + "`" + `
	Name string ` + "`" + `yaml:"name"` + "`" + `
}

type HTTPConfig struct {
	Enabled bool   ` + "`" + `yaml:"enabled"` + "`" + `
	Addr    string ` + "`" + `yaml:"addr"` + "`" + `
}

type GRPCConfig struct {
	Enabled bool   ` + "`" + `yaml:"enabled"` + "`" + `
	Addr    string ` + "`" + `yaml:"addr"` + "`" + `
}

type LogConfig struct {
	Level    string ` + "`" + `yaml:"level"` + "`" + `
	Encoding string ` + "`" + `yaml:"encoding"` + "`" + `
	Output   string ` + "`" + `yaml:"output"` + "`" + `
	Console  *LogConsoleConfig ` + "`" + `yaml:"console"` + "`" + `
	File     *LogFileConfig    ` + "`" + `yaml:"file"` + "`" + `
}

type LogConsoleConfig struct {
	Enabled  *bool  ` + "`" + `yaml:"enabled"` + "`" + `
	Encoding string ` + "`" + `yaml:"encoding"` + "`" + `
	Output   string ` + "`" + `yaml:"output"` + "`" + `
}

type LogFileConfig struct {
	Enabled  *bool  ` + "`" + `yaml:"enabled"` + "`" + `
	Path     string ` + "`" + `yaml:"path"` + "`" + `
	Encoding string ` + "`" + `yaml:"encoding"` + "`" + `
}

type DBConfig struct {
	DSN      string ` + "`" + `yaml:"dsn"` + "`" + `
	MaxOpen  int    ` + "`" + `yaml:"max_open"` + "`" + `
	MaxIdle  int    ` + "`" + `yaml:"max_idle"` + "`" + `
	LogLevel string ` + "`" + `yaml:"log_level"` + "`" + `
}

type RedisConfig struct {
	Addr     string ` + "`" + `yaml:"addr"` + "`" + `
	Password string ` + "`" + `yaml:"password"` + "`" + `
	DB       int    ` + "`" + `yaml:"db"` + "`" + `
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
`

const entityTemplate = `package entity

type Example struct {
	ID   int64  ` + "`" + `json:"id"` + "`" + `
	Name string ` + "`" + `json:"name"` + "`" + `
}
`

const repositoryPortTemplate = `package port

import (
	"context"

	"{{.ServiceName}}/internal/domain/entity"
)

type ExampleRepository interface {
	Create(ctx context.Context, ex *entity.Example) error
	GetByID(ctx context.Context, id int64) (*entity.Example, error)
}
`

const codesTemplate = `package errs

const (
	CodeSuccess       = 0
	CodeInvalidParam = 1
	CodeNotFound     = 2
	CodeInternalError = 9999
)

var CodeMessage = map[int]string{
	CodeSuccess:       "成功",
	CodeInvalidParam: "参数错误",
	CodeNotFound:     "资源不存在",
	CodeInternalError: "系统错误",
}
`

const errorTemplate = `package errs

import "fmt"

type BizError struct {
	Code int
	Msg  string
	Err  error
}

func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e *BizError) Unwrap() error {
	return e.Err
}

func New(code int, msg string, err error) *BizError {
	return &BizError{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}
`

const serviceTemplate = `package service

import (
	"context"

	"{{.ServiceName}}/internal/domain/entity"
	"{{.ServiceName}}/internal/domain/errs"
	"{{.ServiceName}}/internal/domain/port"
)

type ExampleService interface {
	Ping(ctx context.Context) (string, error)
}

type exampleService struct {
	repo port.ExampleRepository
}

func NewExampleService(repo port.ExampleRepository) ExampleService {
	return &exampleService{repo: repo}
}

func (s *exampleService) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
`

const typesTemplate = `package types

type ExampleReq struct {
	Name string ` + "`" + `json:"name" validate:"required"` + "`" + `
}

type ExampleResp struct {
	Message string ` + "`" + `json:"message"` + "`" + `
}
`

const handlerTemplate = `package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ink-yht-code/{{.ServiceName}}/internal/service"
	"github.com/ink-yht-code/sprout"
)

type Handler struct {
	exampleSvc service.ExampleService
}

func NewHandler(exampleSvc service.ExampleService) *Handler {
	return &Handler{exampleSvc: exampleSvc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/ping", sprout.W(h.Ping))
	}
}

func (h *Handler) Ping(ctx *sprout.Context) (sprout.Result, error) {
	msg, err := h.exampleSvc.Ping(ctx)
	if err != nil {
		_ = http.StatusInternalServerError
		return sprout.InternalError(), err
	}
	return sprout.Success("success", map[string]string{"message": msg}), nil
}
`

const wiringTemplate = `package wiring

import (
	"context"

	"github.com/ink-yht-code/{{.ServiceName}}/internal/config"
	"github.com/ink-yht-code/{{.ServiceName}}/internal/service"
	"github.com/ink-yht-code/{{.ServiceName}}/internal/web"
	"github.com/ink-yht-code/sprout/sproutx"
	"gorm.io/gorm"
)

type App struct {
	HTTP *sproutx.Server
	DB   *gorm.DB
}

func InitApp(cfg *config.Config) (*App, error) {
	app := &App{}

	if cfg.DB.DSN != "" {
		db, err := sproutx.NewDB(sproutx.DBConfig{
			DSN:      cfg.DB.DSN,
			MaxOpen:  cfg.DB.MaxOpen,
			MaxIdle:  cfg.DB.MaxIdle,
			LogLevel: cfg.DB.LogLevel,
		})
		if err != nil {
			return nil, err
		}
		app.DB = db
	}

	if cfg.HTTP.Enabled {
		server := sproutx.NewHTTPServer(sproutx.HTTPConfig{
			Enabled: cfg.HTTP.Enabled,
			Addr:    cfg.HTTP.Addr,
		})

		exampleSvc := service.NewExampleService(nil)
		handler := web.NewHandler(exampleSvc)

		handler.RegisterRoutes(server.Engine)

		app.HTTP = server
	}

	return app, nil
}

func (a *App) Run() error {
	if a.HTTP != nil {
		return a.HTTP.Run()
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.HTTP != nil {
		return a.HTTP.Shutdown(ctx)
	}
	return nil
}
`
