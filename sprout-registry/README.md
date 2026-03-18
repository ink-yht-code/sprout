# sprout-registry - ServiceID 注册服务

[![License](https://img.shields.io/badge/License-Proprietary-red.svg)](#license)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![Version](https://img.shields.io/badge/version-v1.0.0-blue.svg)](https://github.com/ink-yht-code/sprout/releases)

Registry 是集中式 ServiceID 分配服务，确保微服务 ID 全局唯一，用于生成标准化的业务错误码。

## 目录

- [为什么需要 Registry？](#为什么需要-registry)
- [特性](#特性)
- [安装](#安装)
- [运行](#运行)
- [配置](#配置)
- [HTTP API](#http-api)
- [客户端使用](#客户端使用)
- [ServiceID 规则](#serviceid-规则)
- [与 sprout-gen 集成](#与-sprout-gen-集成)
- [部署建议](#部署建议)
- [常见问题](#常见问题)
- [许可证](#许可证)

## 版本

当前版本：**v1.0.0**

## 为什么需要 Registry？

在微服务架构中，每个服务需要唯一的标识符来生成业务错误码：

```
业务码 = ServiceID × 10000 + BizCode
```

- **ServiceID** - 服务唯一标识（由 Registry 分配）
- **BizCode** - 业务错误码（服务内部定义）

例如：
- user 服务（ServiceID=101）的用户不存在错误：`1010004`
- order 服务（ServiceID=102）的订单不存在错误：`1020004`

通过 Registry 统一分配 ServiceID，可以：
- 避免服务间 ID 冲突
- 通过错误码快速定位问题服务
- 实现错误码的标准化管理

## 特性

- **幂等分配** - 同名服务多次调用返回相同 ID
- **SQLite 存储** - 轻量级持久化，无需额外依赖
- **HTTP API** - RESTful 接口，易于集成
- **可选认证** - Token 保护，防止未授权访问
- **Go 客户端** - 提供官方 Go 客户端库

## 安装

```bash
# 构建
cd sprout-registry && go build -o bin/registry ./cmd/...

# 或使用 Make
make build-registry
```

## 运行

```bash
# 直接运行
./bin/registry

# 或使用 go run
go run cmd/main.go

# 指定端口和数据文件
go run cmd/main.go -port 8070 -data registry.db
```

默认配置：
- 监听地址: `:18080`
- 数据库: `registry.db`（自动创建）

## 配置

### 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-port` | `18080` | HTTP 服务监听端口 |
| `-data` | `registry.db` | SQLite 数据库文件路径 |
| `-token` | （空） | 认证 Token，设置后启用认证 |

### 环境变量

| 变量 | 说明 |
|------|------|
| `REGISTRY_PORT` | 监听端口 |
| `REGISTRY_DATA` | 数据库文件路径 |
| `REGISTRY_TOKEN` | 认证 Token |

## HTTP API

### 分配 ServiceID

```bash
POST /v1/services:allocate
Content-Type: application/json

{
    "name": "user"
}
```

响应：
```json
{
    "service_id": 101,
    "name": "user"
}
```

**幂等性**：同名服务多次调用返回相同 ID。

### 查询服务

```bash
GET /v1/services/user
```

响应：
```json
{
    "service_id": 101,
    "name": "user"
}
```

**错误响应**：
```json
{
    "error": "service not found"
}
```

### 列出所有服务

```bash
GET /v1/services
```

响应：
```json
{
    "services": [
        {"service_id": 101, "name": "user"},
        {"service_id": 102, "name": "order"},
        {"service_id": 103, "name": "payment"}
    ]
}
```

### 健康检查

```bash
GET /health
```

响应：
```json
{
    "status": "ok"
}
```

### 认证

如果设置了 `-token`，请求需要携带 Token：

```bash
# 启动时设置 Token
go run cmd/main.go -token your-secret-token

# 请求时携带 Token
curl -X POST http://localhost:18080/v1/services:allocate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-secret-token" \
  -d '{"name":"user"}'
```

## 客户端使用

### Go 客户端

```go
import "github.com/ink-yht-code/sprout/sprout-registry/api"

// 创建客户端
client := api.NewClient("http://localhost:18080")

// 可选：设置认证 Token
client.SetToken("your-secret-token")

// 分配 ServiceID
resp, err := client.Allocate(ctx, "user")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("ServiceID: %d\n", resp.ServiceID)

// 查询服务
resp, err := client.Get(ctx, "user")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("ServiceID: %d\n", resp.ServiceID)

// 列出所有服务
services, err := client.List(ctx)
if err != nil {
    log.Fatal(err)
}
for _, svc := range services {
    fmt.Printf("%s: %d\n", svc.Name, svc.ServiceID)
}
```

### HTTP 调用示例

```bash
# 分配 ServiceID
curl -X POST http://localhost:18080/v1/services:allocate \
  -H "Content-Type: application/json" \
  -d '{"name":"user"}'

# 查询服务
curl http://localhost:18080/v1/services/user

# 列出所有服务
curl http://localhost:18080/v1/services

# 健康检查
curl http://localhost:18080/health
```

## ServiceID 规则

### 分配规则

- **起始 ID**: 101（保留 1-100 供系统使用）
- **递增步长**: 1
- **幂等性**: 同名服务返回相同 ID

### 业务码计算

```
业务码 = ServiceID × 10000 + BizCode
```

| BizCode | 含义 |
| ------- | -------- |
| 0 | 成功 |
| 1 | 参数错误 |
| 2 | 未授权 |
| 3 | 无权限 |
| 4 | 未找到 |
| 5 | 冲突 |
| 9999 | 内部错误 |

### 示例

| 服务 | ServiceID | 业务码范围 | 示例 |
| ------- | --------- | --------------- | ------------------- |
| user | 101 | 1010000-1019999 | 用户不存在: 1010004 |
| order | 102 | 1020000-1029999 | 订单不存在: 1020004 |
| payment | 103 | 1030000-1039999 | 支付失败: 1039999 |

### 自定义业务码

服务可以定义自己的业务码：

```go
// internal/domain/errs/codes.go
const (
    CodeUserNotFound    = 4  // 用户不存在
    CodeUserConflict    = 5  // 用户冲突
    CodeInvalidPassword = 10 // 密码错误
    CodeUserBanned      = 11 // 用户被封禁
)

// internal/domain/errs/error.go
var (
    ErrUserNotFound    = New(1010004, "用户不存在")
    ErrUserConflict    = New(1010005, "用户名已存在")
    ErrInvalidPassword = New(1010010, "密码错误")
    ErrUserBanned      = New(1010011, "用户已被封禁")
)
```

## 与 sprout-gen 集成

### 自动分配 ServiceID

`sprout-gen` 创建服务时自动从 Registry 获取 ServiceID：

```bash
# 启动 Registry
cd sprout-registry && go run cmd/main.go

# 创建服务（自动分配 ServiceID）
sprout gen generate user --registry http://localhost:18080
```

生成的 `configs/user.yaml`：

```yaml
service:
  id: 101
  name: user
```

### 手动指定 ServiceID

也可以手动指定 ServiceID：

```bash
sprout gen generate user --service-id 101
```

## 部署建议

### 单实例部署

适合开发环境和小规模生产环境：

```bash
# 直接运行
./registry -port 18080 -data /data/registry.db
```

### Docker 部署

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o registry ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/registry .
EXPOSE 18080
CMD ["./registry"]
```

```bash
# 构建镜像
docker build -t sprout-registry .

# 运行容器
docker run -d \
  --name registry \
  -p 18080:18080 \
  -v /data/registry:/app/data \
  sprout-registry \
  ./registry -data /app/data/registry.db
```

### Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sprout-registry
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sprout-registry
  template:
    metadata:
      labels:
        app: sprout-registry
    spec:
      containers:
      - name: registry
        image: sprout-registry:latest
        ports:
        - containerPort: 18080
        volumeMounts:
        - name: data
          mountPath: /app/data
        args:
        - ./registry
        - -data
        - /app/data/registry.db
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: registry-data
---
apiVersion: v1
kind: Service
metadata:
  name: sprout-registry
spec:
  selector:
    app: sprout-registry
  ports:
  - port: 18080
    targetPort: 18080
```

### 高可用部署

对于大规模生产环境，建议：

1. **使用共享存储** - 将 SQLite 文件放在共享存储（NFS、S3 等）
2. **多实例负载均衡** - 使用 Nginx 或云负载均衡器
3. **数据库迁移** - 迁移到 PostgreSQL 或 MySQL（需要修改代码）

## 常见问题

### Q: Registry 挂了怎么办？

Registry 只在创建新服务时需要，已分配的 ServiceID 会保存在服务配置中。Registry 不可用不影响已运行的服务。

建议：
- 定期备份 `registry.db` 文件
- 监控 Registry 健康状态

### Q: 如何重置 ServiceID？

删除 `registry.db` 文件并重启 Registry。注意：这会导致 ServiceID 重新分配。

### Q: ServiceID 用完了怎么办？

ServiceID 是 int64 类型，理论上不会用完。如果真的接近上限，可以调整分配策略。

### Q: 可以手动指定 ServiceID 吗？

可以，但不推荐。手动指定可能导致冲突。建议通过 Registry 统一分配。

### Q: 如何迁移到其他数据库？

当前只支持 SQLite。如需迁移到 PostgreSQL 或 MySQL，需要修改 `store` 包的实现。

## 许可证

Proprietary License

未经版权所有者书面授权，不得使用、复制、修改或分发本项目的任何部分。

---

Made with ❤️ by [ink-yht-code](https://github.com/ink-yht-code)
