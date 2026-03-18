# sprout-gen - Sprout 代码生成器

[![License](https://img.shields.io/badge/License-Proprietary-red.svg)](#license)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![Version](https://img.shields.io/badge/version-v1.0.0-blue.svg)](https://github.com/ink-yht-code/sprout/releases)

sprout-gen 是 Sprout 框架的代码生成 CLI 工具，支持创建服务骨架、生成 HTTP/gRPC 代码、生成 Repository 等，帮助开发者快速构建符合 DDD + Clean Architecture 的微服务。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [命令详解](#命令详解)
  - [sprout gen new - 创建新项目](#sprout-gen-new---创建新项目)
  - [sprout gen generate - 生成服务骨架](#sprout-gen-generate---生成服务骨架)
  - [sprout gen api - 生成 HTTP 代码](#sprout-gen-api---生成-http-代码)
  - [sprout gen rpc - 生成 gRPC 代码](#sprout-gen-rpc---生成-grpc-代码)
  - [sprout gen repo - 生成 Repository](#sprout-gen-repo---生成-repository)
  - [sprout gen lint - 分层约束检查](#sprout-gen-lint---分层约束检查)
  - [sprout gen entity - 生成实体](#sprout-gen-entity---生成实体)
  - [sprout gen module - 生成模块](#sprout-gen-module---生成模块)
- [.sprout 文件格式](#sprout-文件格式)
- [生成的目录结构](#生成的目录结构)
- [与 Registry 集成](#与-registry-集成)
- [许可证](#许可证)

## 特性

- **服务骨架生成** - 一键创建符合 DDD + Clean Architecture 的服务结构
- **HTTP 代码生成** - 从 .sprout 文件生成 types/handler/routes
- **gRPC 代码生成** - 从 .proto 文件生成 gRPC 代码
- **Repository 生成** - 从 SQL DDL 生成 DAO 和 Repository
- **分层约束检查** - 检查代码分层是否合规
- **增量生成** - 支持增量追加，不覆盖用户代码
- **彩色 CLI 输出** - 美观的命令行帮助信息

## 安装

```bash
go install github.com/ink-yht-code/sprout-gen@latest
```

或从源码构建：

```bash
git clone https://github.com/ink-yht-code/sprout.git
cd sprout
go build -o sprout-gen ./cmd/sprout
```

## 快速开始

### 创建第一个服务

```bash
# 创建新项目
sprout gen new myservice

# 进入项目目录
cd myservice

# 生成服务骨架
sprout gen generate myservice

# 运行服务
go run cmd/main.go
```

服务将在 `http://localhost:8080` 启动。

## 命令详解

### sprout gen new - 创建新项目

创建一个新的 Go 项目目录结构。

```bash
sprout gen new <name>
```

**参数：**
- `name` - 项目名称

**示例：**
```bash
sprout gen new user-service
```

**生成的文件：**
```
user-service/
├── go.mod
├── go.sum
└── cmd/
    └── main.go
```

---

### sprout gen generate - 生成服务骨架

生成完整的服务骨架，包括 DDD 分层结构。

```bash
sprout gen generate <name> [flags]
```

**参数：**
- `name` - 服务名称

**Flags：**
- `--transport` - 传输协议：`http`、`rpc`、`http,rpc`（默认：`http`）
- `--dao` - DAO 类型：`gorm`（默认：`gorm`）
- `--cache` - Cache 类型：`redis`（默认：`redis`）
- `--service-id` - Service ID（0 表示从 Registry 分配）

**示例：**
```bash
# 创建 HTTP 服务
sprout gen generate user --transport http

# 创建 HTTP + gRPC 服务
sprout gen generate order --transport http,rpc

# 创建纯 gRPC 服务
sprout gen generate payment --transport rpc

# 指定 Service ID
sprout gen generate user --service-id 101
```

**生成的目录结构：**
```
user/
├── .sprout              # API 定义文件
├── go.mod               # Go 模块
├── configs/
│   └── user.yaml        # 配置文件
├── cmd/
│   └── main.go          # 入口文件
└── internal/
    ├── config/          # 配置结构
    ├── domain/          # 领域层
    │   ├── entity/      # 实体定义
    │   ├── port/        # Repository 接口
    │   ├── errs/        # 错误定义
    │   └── event/       # 领域事件
    ├── repository/      # 数据访问层
    │   ├── dao/         # DAO 实现
    │   └── cache/       # 缓存层
    ├── service/         # 业务逻辑层
    ├── types/           # 请求/响应类型
    ├── web/             # HTTP 层
    └── wiring/          # 依赖注入
```

---

### sprout gen api - 生成 HTTP 代码

从 `.sprout` 文件生成 HTTP 相关代码。

```bash
sprout gen api <service> [flags]
```

**参数：**
- `service` - 服务名称

**Flags：**
- `-f, --file` - .sprout 文件路径（默认：`<service>/.sprout`）

**示例：**
```bash
# 在服务目录内执行
cd user
sprout gen api user

# 或指定文件路径
sprout gen api user --file user.sprout
```

**生成的文件：**
- `internal/types/types_gen.go` - 类型定义（自动覆盖）
- `internal/web/handlers.go` - Handler 方法骨架（增量追加）
- `internal/web/handler_gen.go` - Handler 结构体和路由（自动覆盖）

**增量生成说明：**
- `types_gen.go` 和 `handler_gen.go` 每次都会重新生成
- `handlers.go` 只会追加新的方法，不会覆盖已有实现

---

### sprout gen rpc - 生成 gRPC 代码

从 `.proto` 文件生成 gRPC 代码。

```bash
sprout gen rpc <service> [flags]
```

**参数：**
- `service` - 服务名称

**Flags：**
- `-p, --proto` - .proto 文件路径（默认：`<service>.proto`）

**示例：**
```bash
# 生成 gRPC 代码
sprout gen rpc user

# 指定 proto 文件
sprout gen rpc user --proto api/user.proto
```

**生成的文件：**
- `api/` - Protoc 生成的 gRPC 代码
- `internal/rpc/server.go` - RPC Server 包装器

**依赖：**
需要安装 `protoc` 和相关插件：
```bash
# 安装 protoc
# macOS: brew install protoc
# Ubuntu: apt install protobuf-compiler

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

---

### sprout gen repo - 生成 Repository

从 SQL DDL 文件生成 DAO 和 Repository 代码。

```bash
sprout gen repo <service> [flags]
```

**参数：**
- `service` - 服务名称

**Flags：**
- `-d, --ddl` - SQL DDL 文件路径
- `-t, --table` - 表名（必填）

**示例：**
```bash
# 从 schema.sql 生成 products 表的 Repository
sprout gen repo product --ddl schema.sql --table products

# 生成 users 表的 Repository
sprout gen repo user --ddl db/schema.sql --table users
```

**生成的文件：**
- `internal/repository/dao/model.go` - DAO 模型（GORM 结构体）
- `internal/repository/dao/dao.go` - DAO 接口和实现
- `internal/repository/repository.go` - Repository 接口和实现

**DDL 文件示例：**
```sql
CREATE TABLE `products` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `stock` int NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**生成的 Model 示例：**
```go
type Product struct {
    Id        int64     `gorm:"column:id"`
    Name      string    `gorm:"column:name"`
    Price     float64   `gorm:"column:price"`
    Stock     int       `gorm:"column:stock"`
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Product) TableName() string {
    return "products"
}
```

---

### sprout gen lint - 分层约束检查

检查代码是否符合 DDD 分层架构规则。

```bash
sprout gen lint <service>
```

**参数：**
- `service` - 服务目录名称

**示例：**
```bash
sprout gen lint user
```

**检查规则：**

| 层级 | 允许依赖 | 禁止依赖 |
|------|---------|---------|
| `web` | `service`, `types`, `domain` | `repository`, `dao` |
| `service` | `repository`, `domain` | `web`, `dao` |
| `repository` | `dao`, `domain` | `web`, `service` |
| `dao` | `domain` | `web`, `service`, `repository` |
| `domain` | 无 | 所有其他层 |

**输出示例：**
```
检查分层约束: user

✓ web 层检查通过
✓ service 层检查通过
✓ repository 层检查通过
✓ dao 层检查通过
✓ domain 层检查通过

分层约束检查完成，无违规
```

---

### sprout gen entity - 生成实体

生成领域实体和对应的 Repository 接口。

```bash
sprout gen entity <name>
```

**参数：**
- `name` - 实体名称

**示例：**
```bash
sprout gen entity User
```

**生成的文件：**
- `internal/domain/entity/user.go` - 实体定义
- `internal/domain/port/user_repository.go` - Repository 接口

---

### sprout gen module - 生成模块

生成完整的业务模块（包含 entity + service + repository）。

```bash
sprout gen module <name>
```

**参数：**
- `name` - 模块名称

**示例：**
```bash
sprout gen module Order
```

**生成的文件：**
```
internal/
├── domain/
│   ├── entity/order.go
│   └── port/order_repository.go
├── repository/
│   └── order.go
└── service/
    └── order.go
```

---

## .sprout 文件格式

`.sprout` 文件定义 HTTP API 的类型和服务。

### 基本语法

```go
syntax = "v1"

// 类型定义
type <TypeName> {
    <FieldName> <FieldType> `json:"<json_name>" validate:"<rules>"`
}

// 服务器配置
server {
    prefix "<url_prefix>"
}

// 服务定义
service <ServiceName> {
    public {
        <METHOD> "<path>" <HandlerName>(<ReqType>) -> <RespType>
    }
    
    private {
        <METHOD> "<path>" <HandlerName>(<ReqType>) -> <RespType>
    }
}
```

### 字段类型

| 类型 | Go 类型 | 说明 |
|------|---------|------|
| `string` | `string` | 字符串 |
| `int` | `int` | 整数 |
| `int64` | `int64` | 64位整数 |
| `float64` | `float64` | 浮点数 |
| `bool` | `bool` | 布尔值 |
| `time.Time` | `time.Time` | 时间 |

### HTTP 方法

支持：`GET`, `POST`, `PUT`, `DELETE`, `PATCH`

### 完整示例

```go
syntax = "v1"

// 请求类型
type CreateUserReq {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=0,max=150"`
}

// 响应类型
type CreateUserResp {
    Id    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type GetUserResp {
    Id    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UpdateUserReq {
    Id    int64  `json:"id" validate:"required"`
    Name  string `json:"name"`
    Email string `json:"email" validate:"email"`
}

type ListUsersResp {
    Items []UserItem `json:"items"`
    Total int64      `json:"total"`
}

type UserItem {
    Id    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// 服务器配置
server {
    prefix "/api/v1"
}

// 服务定义
service user {
    // 公开接口（不需要认证）
    public {
        POST "/users" CreateUser(CreateUserReq) -> CreateUserResp
        GET "/users" ListUsers -> ListUsersResp
        GET "/users/:id" GetUser -> GetUserResp
    }
    
    // 私有接口（需要认证）
    private {
        PUT "/users/:id" UpdateUser(UpdateUserReq) -> GetUserResp
        DELETE "/users/:id" DeleteUser
    }
}
```

---

## 生成的目录结构

完整的 DDD 分层目录结构：

```
<service>/
├── .sprout                    # API 定义文件（手动编辑）
├── go.mod                     # Go 模块文件
├── go.sum                     # 依赖锁定
├── README.md                  # 服务说明
├── configs/
│   └── <service>.yaml         # 配置文件
├── cmd/
│   └── main.go                # 程序入口
└── internal/
    ├── config/
    │   └── config.go          # 配置结构
    ├── domain/                # 领域层
    │   ├── entity/            # 实体定义
    │   ├── port/              # Repository 接口
    │   ├── errs/              # 错误定义
    │   │   ├── codes.go       # 错误码
    │   │   └── error.go       # 业务错误
    │   └── event/             # 领域事件
    ├── repository/            # 数据访问层
    │   ├── dao/               # DAO 实现
    │   │   ├── model.go       # GORM 模型
    │   │   └── dao.go         # DAO 接口
    │   ├── cache/             # 缓存层
    │   └── <service>.go       # Repository 实现
    ├── service/               # 业务逻辑层
    │   └── <service>.go       # Service 实现
    ├── types/                 # 请求/响应类型
    │   └── types_gen.go       # 类型定义（生成）
    ├── web/                   # HTTP 层
    │   ├── handler_gen.go     # Handler 骨架（生成）
    │   └── handlers.go        # Handler 实现
    ├── event/                 # 事件层
    │   └── consumer.go        # 事件消费者
    └── wiring/                # 依赖注入
        ├── wiring.go          # Wire 定义
        └── wire_gen.go        # Wire 生成
```

---

## 与 Registry 集成

sprout-registry 是服务注册中心，用于统一分配 ServiceID。

### 启动 Registry

```bash
cd sprout-registry
go run cmd/main.go

# Registry 默认监听 :18080
```

### 使用 Registry 创建服务

```bash
# 自动从 Registry 获取 ServiceID
sprout gen generate user --registry http://localhost:18080
```

### ServiceID 规范

- ServiceID 用于生成唯一错误码
- 错误码 = ServiceID × 10000 + BizCode
- 例如：user 服务（ServiceID=101）
  - 参数错误：1010001
  - 未授权：1010002
  - 内部错误：1019999

---

## 许可证

Proprietary License

未经版权所有者书面授权，不得使用、复制、修改或分发本项目的任何部分。

---

Made with ❤️ by [ink-yht-code](https://github.com/ink-yht-code)
