package scaffoldtmpl

// ReadmeTmpl 用于生成项目 README.md 的模板
var ReadmeTmpl = `# {{.Name}}

{{.Name}} 服务

## 快速开始

### 1. 安装依赖

` + "`" + `bash
go mod tidy
` + "`" + `

**注意**：如果 ` + "`" + `go mod tidy` + "`" + ` 报错，请确保已配置 GOPRIVATE 环境变量。项目已生成 ` + "`" + `.env` + "`" + ` 文件，可以使用以下命令加载：
` + "`" + `bash
# Linux/Mac
source .env

# Windows PowerShell
Get-Content .env | ForEach-Object { $var = $_.Split('='); [System.Environment]::SetEnvironmentVariable($var[0], $var[1]) }
` + "`" + `

或者手动设置环境变量：

` + "`" + `bash
export GOPRIVATE=*.ink-yht-code
` + "`" + `

### 2. 配置

编辑 ` + "`" + `configs/config.yaml` + "`" + `:

` + "`" + `yaml
service:
  id: {{.ServiceID}}
  name: {{.Name}}

http:
  enabled: true
  addr: ":8080"

jwt:
  secret: "your-secret-key-change-in-production"

db:
  dsn: "user:pass@tcp(127.0.0.1:3306)/{{.Name}}?charset=utf8mb4&parseTime=True&loc=Local"
` + "`" + `

### 3. 运行

` + "`" + `bash
go run cmd/main.go
` + "`" + `

## 目录结构

` + "`" + `
{{.Name}}/
├── cmd/main.go              # 入口
├── configs/config.yaml       # 配置
├── internal/
│   ├── config/              # 配置解析
│   ├── domain/
│   │   ├── entity/          # 实体
│   │   ├── port/            # 仓储接口
│   │   └── errs/            # 错误码
│   ├── repository/
│   │   └── dao/             # DAO 实现
│   ├── service/             # 业务逻辑
│   ├── types/               # 请求/响应类型
│   ├── web/                 # HTTP Handler
│   └── wiring/              # 依赖注入
├── {{.Name}}.gint           # API 定义
└── go.mod
` + "`" + `

## API

查看 ` + "`" + `{{.Name}}.gint` + "`" + ` 文件了解 API 定义。

### 健康检查

` + "`" + `bash
curl http://localhost:8080/health
` + "`" + `

## 开发

### 重新生成代码

` + "`" + `bash
gint-gen api {{.Name}}
` + "`" + `

### 运行测试

` + "`" + `bash
go test ./...
` + "`" + `
`
