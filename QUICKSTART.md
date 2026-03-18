# Sprout 快速开始

本指南将帮助你在 5 分钟内上手 Sprout 框架，创建第一个微服务。

## 前置要求

- Go 1.25 或更高版本
- MySQL 5.7+ 或 8.0+
- Redis 5.0+

## 安装

### 1. 安装 sprout 和 sproutx

```bash
go get github.com/ink-yht-code/sprout
go get github.com/ink-yht-code/sproutx
```

### 2. 安装 sprout-gen

```bash
go install github.com/ink-yht-code/sprout-gen@latest
```

### 3. 启动 sprout-registry（可选）

```bash
cd sprout-registry
go run cmd/main.go
```

Registry 默认监听在 `:18080`。

## 创建第一个服务

### 步骤 1：创建服务骨架

```bash
# 创建一个名为 user 的 HTTP 服务
sprout-gen new service user --transport http

# 进入服务目录
cd user
```

这会创建以下目录结构：

```
user/
├── user.sprout              # API 定义文件
├── go.mod                   # Go 模块
├── README.md                # 服务说明
├── configs/
│   └── user.yaml           # 服务配置
├── cmd/
│   └── main.go             # 入口文件
└── internal/
    ├── config/
    │   └── config.go       # 配置结构
    ├── domain/             # 领域层
    ├── repository/         # 数据访问层
    ├── service/            # 业务逻辑层
    ├── types/              # 请求/响应结构
    ├── web/                # HTTP 层
    └── wiring/            # 依赖组装
```

### 步骤 2：编辑 API 定义

打开 `user.sprout` 文件，定义你的 API：

```go
syntax = "v1"

type CreateUserReq {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

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

server {
    prefix "/api"
}

service user {
    public {
        POST "/users" CreateUser(CreateUserReq) -> CreateUserResp
        GET "/users/:id" GetUser -> GetUserResp
    }
}
```

### 步骤 3：生成代码

```bash
sprout-gen api user
```

这会生成：
- `internal/types/types_gen.go` - 类型定义
- `internal/web/handler_gen.go` - Handler 骨架
- `internal/web/user_handlers.go` - Handler 实现（增量追加）

### 步骤 4：实现业务逻辑

编辑 `internal/service/user.go`：

```go
package service

import (
    "context"
    "github.com/your-project/user-service/internal/domain"
    "github.com/your-project/user-service/internal/types"
)

type UserService interface {
    CreateUser(ctx context.Context, req *types.CreateUserReq) (*domain.User, error)
    GetUser(ctx context.Context, id int64) (*domain.User, error)
}
```

编辑 `internal/web/user_handlers.go`：

```go
package web

import (
    "github.com/ink-yht-code/sprout"
    "github.com/ink-yht-code/sprout/context"
    "github.com/your-project/user-service/internal/types"
)

func (h *Handler) CreateUser(ctx *context.Context, req *types.CreateUserReq) (sprout.Result, error) {
    user, err := h.svc.CreateUser(ctx, req)
    if err != nil {
        return sprout.Error("创建用户失败"), err
    }
    return sprout.Success("创建成功", user), nil
}

func (h *Handler) GetUser(ctx *context.Context) (sprout.Result, error) {
    id := ctx.Param("id").Int64()
    user, err := h.svc.GetUser(ctx, id)
    if err != nil {
        return sprout.NotFound("用户不存在"), err
    }
    return sprout.Success("查询成功", user), nil
}
```

### 步骤 5：配置服务

编辑 `configs/user.yaml`：

```yaml
service:
  id: 101
  name: user

http:
  enabled: true
  addr: ":8080"

log:
  level: "info"
  encoding: "json"
  output: "stdout"

db:
  dsn: "user:pass@tcp(127.0.0.1:3306)/user_db?charset=utf8mb4&parseTime=True&loc=Local"
  max_open: 100
  max_idle: 10
  log_level: "info"

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0
```

### 步骤 6：运行服务

```bash
go run cmd/main.go
```

服务启动后，你可以访问：

```bash
# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# 查询用户
curl http://localhost:8080/api/users/1
```

## 下一步

- 阅读 [sprout 文档](sprout/README.md) 了解 HTTP 框架的更多功能
- 阅读 [sproutx 文档](sproutx/README.md) 了解运行时框架的更多功能
- 阅读 [sprout-gen 文档](sprout-gen/README.md) 了解代码生成器的更多功能
- 查看 [examples](examples/) 目录中的示例项目

## 常见问题

### Q: 如何添加新的 API？

A: 编辑 `user.sprout` 文件，添加新的 API 定义，然后运行 `sprout-gen api user` 重新生成代码。

### Q: 如何使用 Session？

A: 使用 `sprout.S` 或 `sprout.BS` 包装器，它们会自动处理 Session。

```go
r.GET("/profile", sprout.S(func(ctx *sprout.Context, sess sprout.Session) (sprout.Result, error) {
    userID := sess.Claims().UserId
    return sprout.Success("查询成功", userID), nil
}))
```

### Q: 如何使用事务？

A: 使用 `sproutx.TxManager` 在 Service 层开启事务。

```go
err := txMgr.Do(ctx, func(ctx context.Context) error {
    repo.Save(ctx, user)
    repo.Save(ctx, profile)
    return nil
})
```

### Q: 如何使用 JWT？

A: 使用 `sprout.NewJWTManager` 创建 JWT 管理器。

```go
manager := sprout.NewJWTManager(sprout.NewJWTOptions("your-secret-key", 7200, 604800))
claims := jwt.Claims{
    UserId: "12345",
    SSID:   "session-id",
}
tokenPair, err := manager.GenerateTokenPair(claims)
```

## 获取帮助

- 查看 [完整文档](README.md)
- 提交 [Issue](https://github.com/ink-yht-code/sprout/issues)
- 加入 [讨论区](https://github.com/ink-yht-code/sprout/discussions)

---

开始使用 Sprout，构建你的微服务吧！🚀
