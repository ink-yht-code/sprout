# Sprout - Go 微服务开发框架

[![License](https://img.shields.io/badge/License-Proprietary-red.svg)](#license)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![Version](https://img.shields.io/badge/version-v1.0.0-blue.svg)](https://github.com/ink-yht-code/sprout/releases)

Sprout 是一个完整的 Go 微服务开发框架，提供从代码生成到运行时的全栈解决方案，帮助开发者快速构建符合 DDD + Clean Architecture 的微服务。

## 🌱 项目简介

Sprout 整合了四个核心模块，提供完整的微服务开发体验：

- **sprout** - HTTP 框架库，提供 Handler 包装器、中间件、Session 管理、JWT 认证、权限控制等功能
- **sproutx** - 运行时框架，提供日志、数据库、Redis、事务、HTTP/gRPC 服务器、健康检查等基础设施
- **sprout-gen** - 代码生成器，快速创建服务骨架、生成 HTTP/gRPC 代码、Repository 等
- **sprout-registry** - 服务注册中心，统一分配 ServiceID，确保微服务 ID 全局唯一

## ✨ 特性

### 🚀 快速开发
- 一键创建符合 DDD + Clean Architecture 的服务结构
- 自动生成类型定义、Handler 骨架和路由注册
- 支持增量生成，不覆盖用户代码

### 🛡️ 开箱即用
- 丰富的中间件（CORS、CSRF、限流、活跃连接限制、访问日志等）
- Session 管理（Memory 和 Redis 存储，支持 Cookie 和 Header 传递）
- JWT 认证（支持双 Token 机制）
- Casbin 权限管理（基于 RBAC）

### 🏗️ 生产级基础设施
- 结构化日志（基于 zap，支持链路追踪）
- 事务管理（基于 ctx 的事务传递，DAO 层无感知）
- 数据库初始化（GORM，支持连接池配置）
- Redis 初始化（基于 go-redis/v9）
- HTTP/gRPC 服务器（内置中间件，支持优雅关闭）
- 健康检查（HTTP 和 gRPC，支持依赖检查）

### 📦 统一管理
- ServiceID 统一分配，避免冲突
- 标准化错误码（ServiceID × 10000 + BizCode）
- 服务注册与发现

## 📦 安装

### 安装 sprout 和 sproutx

```bash
go get github.com/ink-yht-code/sprout
go get github.com/ink-yht-code/sproutx
```

### 安装 sprout-gen

```bash
go install github.com/ink-yht-code/sprout-gen@latest
```

### 运行 sprout-registry

```bash
cd sprout-registry
go run cmd/main.go
```

## 🚀 快速开始

### 1. 创建新服务

```bash
# 创建 HTTP 服务
sprout-gen new service user --transport http

# 进入服务目录
cd user
```

### 2. 编辑 API 定义

创建 `user.sprout` 文件：

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

### 3. 生成代码

```bash
sprout-gen api user
```

### 4. 实现业务逻辑

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
```

### 5. 运行服务

```bash
go run cmd/main.go
```

## 📚 文档

- [快速开始](QUICKSTART.md) - 5 分钟上手 Sprout
- [sprout 文档](sprout/README.md) - HTTP 框架库详细文档
- [sproutx 文档](sproutx/README.md) - 运行时框架详细文档
- [sprout-gen 文档](sprout-gen/README.md) - 代码生成器详细文档
- [sprout-registry 文档](sprout-registry/README.md) - 服务注册中心详细文档

## 🏗️ 项目结构

```
sprout/                          # 项目根目录
├── sprout/                      # HTTP 框架库
│   ├── core/                    # 核心功能
│   ├── context/                 # Context 增强
│   ├── middlewares/             # 中间件
│   ├── session/                 # Session 管理
│   ├── jwt/                     # JWT 管理
│   └── casbin/                  # 权限管理
├── sproutx/                     # 运行时框架
│   ├── app/                     # 应用启动器
│   ├── log/                     # 结构化日志
│   ├── db/                      # 数据库初始化
│   ├── redis/                   # Redis 初始化
│   ├── tx/                      # 事务管理
│   ├── httpx/                   # HTTP 服务器
│   ├── rpc/                     # gRPC 服务器
│   └── health/                  # 健康检查
├── sprout-gen/                  # 代码生成器
│   ├── cmd/                     # CLI 命令
│   ├── parser/                  # .sprout 文件解析
│   ├── generator/               # 代码生成器
│   ├── registry/                # Registry 客户端
│   └── template/                # 代码模板
├── sprout-registry/             # 服务注册中心
│   ├── cmd/                     # 服务入口
│   ├── api/                     # HTTP API
│   └── store/                   # 存储抽象
├── examples/                    # 示例项目
├── docs/                        # 文档
├── scripts/                     # 脚本
├── go.work                      # Go workspace
├── go.mod                       # 根模块
└── README.md                    # 项目说明
```

## 🔧 技术栈

| 类别 | 技术选型 | 版本 | 用途 |
|------|---------|------|------|
| Web框架 | Gin | v1.12.0 | HTTP 服务 |
| ORM | GORM | v1.25.12 | 数据库操作 |
| 日志 | Zap | v1.27.0 | 结构化日志 |
| Redis | go-redis/v9 | v9.2.1 | 缓存、Session |
| gRPC | google.golang.org/grpc | v1.79.1 | RPC 通信 |
| JWT | golang-jwt/jwt/v5 | v5.0.0 | JWT 认证 |
| Casbin | casbin/v2 | v2.134.0 | 权限控制 |
| Cobra | spf13/cobra | v1.8.0 | CLI 工具 |

## 📝 错误码规范

业务码 = ServiceID × 10000 + BizCode

| BizCode | 含义 |
| ------- | -------- |
| 0 | 成功 |
| 1 | 参数错误 |
| 2 | 未授权 |
| 3 | 无权限 |
| 4 | 未找到 |
| 5 | 冲突 |
| 9999 | 内部错误 |

例如 user 服务（ServiceID=101）：
- 参数错误：1010001
- 未授权：1010002
- 内部错误：1019999

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

Proprietary License

未经版权所有者书面授权，不得使用、复制、修改或分发本项目的任何部分。

## 📮 联系方式

- 项目主页: [https://github.com/ink-yht-code/sprout](https://github.com/ink-yht-code/sprout)
- 问题反馈: [https://github.com/ink-yht-code/sprout/issues](https://github.com/ink-yht-code/sprout/issues)

---

Made with ❤️ by [ink-yht-code](https://github.com/ink-yht-code)
