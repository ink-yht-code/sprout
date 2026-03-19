package scaffoldtmpl

// This file is migrated from legacy gint-gen/template/templates.go.

// GoModTmpl go.mod 模板
var GoModTmpl = `module {{.Name}}

go 1.25

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/ink-yht-code/sprout {{.SproutVersion}}
	github.com/ink-yht-code/sprout/sproutx {{.SproutxVersion}}
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
	gopkg.in/yaml.v3 v3.0.1
)
`

// GintTmpl .gint 文件模板
var GintTmpl = `syntax = "v1"

// 请求类型
type HelloReq {
    Name string ` + "`" + `json:"name" validate:"required"` + "`" + `
}

// 响应类型
type HelloResp {
    Message string ` + "`" + `json:"message"` + "`" + `
}

// 服务器配置
server {
    prefix "/api"
}

// 服务定义
service {{.Name}} {
    // 公开接口（不需要认证）
    public {
        GET "/hello" Hello(HelloReq) -> HelloResp
    }

    // 私有接口（需要认证）
    private {
        // TODO: 添加需要认证的接口
    }
}
`

// ConfigYamlTmpl 配置文件模板
var ConfigYamlTmpl = `# {{.Name}} service config
# 注意：此文件由 gint-gen 生成，你可以按需修改。

# 服务基础信息
service:
  id: {{.ServiceID}} # ServiceID：由 registry 分配（或你自行维护），用于错误码分段等
  name: {{.Name}}

# HTTP 服务配置
http:
  enabled: {{.HasHTTP}} # 是否启用 HTTP
  addr: ":8080" # 监听地址

# gRPC 服务配置
grpc:
  enabled: {{.HasRPC}} # 是否启用 gRPC
  addr: ":9090" # 监听地址

# JWT 配置
jwt:
  secret: "your-secret-key-change-in-production" # 生产环境务必替换
  access_expire: "2h" # AccessToken 过期时间
  refresh_expire: "168h" # RefreshToken 过期时间
  issuer: "{{.Name}}" # 签发者

# 数据库配置（MySQL DSN 示例）
# 如果你本地没有数据库，可保留示例；框架会在连接失败时降级为内存 DAO（demo 可正常启动）。
db:
  dsn: "user:pass@tcp(127.0.0.1:3306)/{{.Name}}?charset=utf8mb4&parseTime=True&loc=Local"
  max_open: 100 # 最大连接数
  max_idle: 10  # 最大空闲连接数
  log_level: "info" # gorm 日志级别: silent|error|warn|info

# Redis 配置
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

# 日志配置
log:
  # 日志级别: debug|info|warn|error
  level: "info"
  # 旧版兼容配置: json|console
  # - output 为 stdout/stderr 时：按 encoding 输出到终端
  # - output 为文件路径时：默认启用 "终端 console + 文件 json" 的双输出
  encoding: "json"
  # 旧版兼容配置: stdout|stderr|文件路径
  output: "stdout"

  # 推荐：终端 console（可读）+ 文件 json（便于采集入库）
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

# Outbox 配置（事件投递）
outbox:
  enabled: true
  batch_size: 100
  poll_interval: "5s"
`

// MainTmpl main.go 模板
var MainTmpl = `package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.Name}}/internal/wiring"
)

func main() {
	// 加载配置
	cfg, err := wiring.LoadConfig("configs/{{.Name}}.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 构建应用
	app, err := wiring.BuildApp(cfg)
	if err != nil {
		fmt.Printf("Failed to build app: %v\n", err)
		os.Exit(1)
	}

	// 启动
	go func() {
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Printf("Service {{.Name}} started\n")

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Shutdown(ctx)
	fmt.Println("Service exited")
}
`

// NOTE: For brevity, remaining legacy templates are not yet included.
// TODO: If you want full parity, we will continue migrating the rest of templates from legacy file.
