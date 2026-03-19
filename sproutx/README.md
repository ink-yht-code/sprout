# sproutx - 运行时框架

[![License](https://img.shields.io/badge/License-Proprietary-red.svg)](#license)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![Version](https://img.shields.io/badge/version-v1.0.0-blue.svg)](https://github.com/ink-yht-code/sprout/releases)

sproutx 提供微服务运行时基础设施组件，包括日志、数据库、Redis、HTTP/gRPC 服务器、事务管理、健康检查等，帮助开发者快速构建生产级的微服务。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [核心模块](#核心模块)
  - [app - 应用启动器](#app---应用启动器)
  - [log - 结构化日志](#log---结构化日志)
  - [db - 数据库初始化](#db---数据库初始化)
  - [redis - Redis 初始化](#redis---redis-初始化)
  - [tx - 事务管理](#tx---事务管理)
  - [httpx - HTTP 服务器](#httpx---http-服务器)
  - [rpc - gRPC 服务器](#rpc---grpc-服务器)
  - [health - 健康检查](#health---健康检查)
  - [outbox - 事务消息](#outbox---事务消息)
- [配置说明](#配置说明)
- [许可证](#许可证)

## 特性

- **结构化日志** - 基于 zap，支持 ctx 注入 request_id，便于链路追踪
- **事务管理** - 基于 ctx 的事务传递，简化事务处理，DAO 层无感知
- **数据库初始化** - GORM 初始化，支持连接池配置和日志级别设置
- **Redis 初始化** - Redis 客户端初始化，基于 go-redis/v9
- **HTTP 服务器** - Gin 服务器，内置常用中间件，支持优雅关闭
- **gRPC 服务器** - gRPC 服务器，内置拦截器，支持优雅关闭
- **健康检查** - HTTP 和 gRPC 健康检查，支持依赖检查
- **优雅关闭** - 支持信号监听和超时关闭
- **应用启动器** - 统一初始化所有组件，简化应用启动流程
- **事务消息** - Outbox 模式，保证消息可靠发送

## 安装

```bash
go get github.com/ink-yht-code/sprout
```

## 快速开始

```go
package main

import (
    "github.com/ink-yht-code/sprout/sproutx"
)

func main() {
    application, err := sproutx.NewApp(&sproutx.Config{
        Service: sproutx.ServiceConfig{ID: 101, Name: "user"},
        HTTP:    sproutx.HTTPConfig{Enabled: true, Addr: ":8080"},
        Log:     sproutx.LogConfig{Level: "info", Encoding: "json"},
        DB:      sproutx.DBConfig{DSN: "user:pass@tcp(127.0.0.1:3306)/user_db"},
        Redis:   sproutx.RedisConfig{Addr: "127.0.0.1:6379"},
    })
    if err != nil {
        panic(err)
    }
    
    // 启动应用
    application.Run()
}
```

## 核心模块

### app - 应用启动器

统一初始化和管理所有运行时组件。

#### 创建应用

```go
import "github.com/ink-yht-code/sprout/sproutx"

application, err := sproutx.NewApp(&sproutx.Config{
    Service: sproutx.ServiceConfig{
        ID:   101,
        Name: "user",
    },
    HTTP: sproutx.HTTPConfig{
        Enabled: true,
        Addr:    ":8080",
    },
    GRPC: sproutx.GRPCConfig{
        Enabled: true,
        Addr:    ":9090",
    },
    Log: sproutx.LogConfig{
        Level:    "info",
        Encoding: "json",
        Output:   "stdout",
    },
    DB: sproutx.DBConfig{
        DSN:      "user:pass@tcp(127.0.0.1:3306)/user_db",
        MaxOpen:  100,
        MaxIdle:  10,
        LogLevel: "info",
    },
    Redis: sproutx.RedisConfig{
        Addr:     "127.0.0.1:6379",
        Password: "",
        DB:       0,
    },
})
```

#### 应用生命周期

```go
// 启动应用（阻塞）
application.Run()

// 或异步启动
go application.Run()

// 等待关闭信号
application.WaitForShutdown()

// 主动关闭
application.Shutdown()
```

#### 注册初始化钩子

```go
application.RegisterInit(func(app *app.App) error {
    // 自定义初始化逻辑
    return nil
})
```

---

### log - 结构化日志

基于 zap 的结构化日志，支持链路追踪。

#### 初始化

```go
import "github.com/ink-yht-code/sprout/sproutx/log"

err := log.Init(log.Config{
    Level:    "info",     // debug, info, warn, error
    Encoding: "json",     // json, console
    Output:   "stdout",   // stdout, file path
    File: log.FileConfig{
        Path:     "/var/log/app.log",
        MaxSize:  100,  // MB
        MaxAge:   30,   // days
        Compress: true,
    },
})
```

#### 使用日志

```go
// 全局日志
log.L().Info("service started",
    zap.String("service", "user"),
    zap.Int("port", 8080),
)

// 带 Context 的日志（自动注入 request_id）
log.Ctx(ctx).Info("handling request",
    zap.String("path", "/users"),
)

// 不同级别
log.L().Debug("debug message")
log.L().Info("info message")
log.L().Warn("warn message")
log.L().Error("error message")

// 带字段的日志
log.L().Error("database error",
    zap.Error(err),
    zap.String("operation", "insert"),
)

// Sugar 风格（性能稍低，更简洁）
log.L().Sugar().Infow("user created",
    "user_id", 123,
    "name", "John",
)
```

#### 链路追踪

日志自动从 Context 提取：
- `request_id` - 请求唯一标识
- `trace_id` - 分布式追踪 ID
- `span_id` - 当前 Span ID

```go
// 在中间件中注入
func TraceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        ctx = log.WithRequestID(c.Request.Context(), requestID)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

---

### db - 数据库初始化

GORM 数据库初始化，支持连接池和日志配置。

#### 初始化

```go
import "github.com/ink-yht-code/sprout/sproutx/db"

database, err := db.New(db.Config{
    DSN:      "user:pass@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
    MaxOpen:  100,        // 最大打开连接数
    MaxIdle:  10,         // 最大空闲连接数
    LogLevel: "info",     // silent, error, warn, info
})
if err != nil {
    panic(err)
}
```

#### 使用数据库

```go
// 自动迁移
err := database.AutoMigrate(&User{}, &Profile{})

// 查询
var user User
err := database.WithContext(ctx).First(&user, 1).Error

// 创建
err := database.WithContext(ctx).Create(&user).Error

// 更新
err := database.WithContext(ctx).Model(&user).Update("name", "John").Error

// 删除
err := database.WithContext(ctx).Delete(&user).Error
```

#### 连接池监控

```go
stats := database.Stats()
fmt.Printf("OpenConnections: %d\n", stats.OpenConnections)
fmt.Printf("InUse: %d\n", stats.InUse)
fmt.Printf("Idle: %d\n", stats.Idle)
```

---

### redis - Redis 初始化

Redis 客户端初始化，基于 go-redis/v9。

#### 初始化

```go
import "github.com/ink-yht-code/sprout/sproutx/redis"

client, err := redis.New(redis.Config{
    Addr:         "127.0.0.1:6379",
    Password:     "",
    DB:           0,
    PoolSize:     100,
    MinIdleConns: 10,
})
if err != nil {
    panic(err)
}
```

#### 使用 Redis

```go
import "github.com/redis/go-redis/v9"

// 设置值
err := client.Set(ctx, "key", "value", time.Hour).Err()

// 获取值
val, err := client.Get(ctx, "key").Result()

// 设置过期
err := client.Expire(ctx, "key", time.Hour).Err()

// 删除
err := client.Del(ctx, "key").Err()

// Hash
err := client.HSet(ctx, "user:1", "name", "John").Err()
name, err := client.HGet(ctx, "user:1", "name").Result()

// List
err := client.LPush(ctx, "queue", "item1", "item2").Err()
items, err := client.LRange(ctx, "queue", 0, -1).Result()

// Set
err := client.SAdd(ctx, "set", "member1", "member2").Err()
members, err := client.SMembers(ctx, "set").Result()

// Sorted Set
err := client.ZAdd(ctx, "rank", redis.Z{Score: 100, Member: "user1"}).Err()

// 分布式锁
locked, err := client.SetNX(ctx, "lock:key", "1", time.Second*30).Result()
```

---

### tx - 事务管理

基于 Context 的事务管理，DAO 层无感知。

#### 创建事务管理器

```go
import "github.com/ink-yht-code/sprout/sproutx/tx"

txMgr := tx.NewTxManager(db)
```

#### 使用事务

```go
// 执行事务
err := txMgr.Do(ctx, func(ctx context.Context) error {
    // 在同一个事务中执行多个操作
    if err := userRepo.Save(ctx, user); err != nil {
        return err  // 自动回滚
    }
    if err := profileRepo.Save(ctx, profile); err != nil {
        return err  // 自动回滚
    }
    return nil  // 自动提交
})
```

#### DAO 层实现

DAO 层不需要知道事务的存在，只需要使用传入的 Context：

```go
type UserDAO struct {
    db *gorm.DB
}

func (d *UserDAO) Save(ctx context.Context, user *User) error {
    // 自动使用 Context 中的事务（如果有）
    return d.db.WithContext(ctx).Create(user).Error
}

func (d *UserDAO) FindByID(ctx context.Context, id int64) (*User, error) {
    var user User
    err := d.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}
```

#### 嵌套事务

```go
err := txMgr.Do(ctx, func(ctx context.Context) error {
    // 外层事务
    
    // 内层操作会加入同一个事务
    if err := txMgr.Do(ctx, func(ctx context.Context) error {
        return repoA.Save(ctx, a)
    }); err != nil {
        return err
    }
    
    return repoB.Save(ctx, b)
})
```

---

### httpx - HTTP 服务器

基于 Gin 的 HTTP 服务器，内置中间件和优雅关闭。

#### 创建服务器

```go
import "github.com/ink-yht-code/sprout/sproutx/httpx"

server, err := httpx.NewServer(httpx.Config{
    Enabled: true,
    Addr:    ":8080",
    Mode:    "release",  // debug, release, test
})
if err != nil {
    panic(err)
}
```

#### 注册路由

```go
// 获取 Gin Engine
engine := server.Engine()

// 注册中间件
engine.Use(gin.Recovery())
engine.Use(gin.Logger())

// 注册路由
engine.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
})

api := engine.Group("/api/v1")
{
    api.POST("/users", userHandler.Create)
    api.GET("/users/:id", userHandler.GetByID)
}
```

#### 启动服务器

```go
// 阻塞启动
server.Run()

// 或异步启动
go server.Run()

// 优雅关闭
server.Shutdown(context.Background())
```

#### 内置中间件

```go
// CORS
server.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))

// 限流
server.Use(rateLimit.New(rateLimit.Config{
    Rate:  100,
    Burst: 200,
}))

// 访问日志
server.Use(accesslog.New(accesslog.Config{
    SkipPaths: []string{"/health", "/metrics"},
}))
```

---

### rpc - gRPC 服务器

gRPC 服务器，内置拦截器和优雅关闭。

#### 创建服务器

```go
import "github.com/ink-yht-code/sprout/sproutx/rpc"

server, err := rpc.NewServer(rpc.Config{
    Enabled: true,
    Addr:    ":9090",
})
if err != nil {
    panic(err)
}
```

#### 注册服务

```go
import pb "github.com/your-org/api/proto"

// 注册 gRPC 服务
pb.RegisterUserServiceServer(server.Server(), userServer)

// 注册健康检查
health.Register(server.Server())
```

#### 启动服务器

```go
// 阻塞启动
server.Run()

// 或异步启动
go server.Run()

// 优雅关闭
server.Shutdown(context.Background())
```

#### 内置拦截器

```go
// 日志拦截器
server.AddUnaryInterceptor(logging.UnaryServerInterceptor())

// 恢复拦截器
server.AddUnaryInterceptor(recovery.UnaryServerInterceptor())

// 链路追踪拦截器
server.AddUnaryInterceptor(otel.UnaryServerInterceptor())
```

---

### health - 健康检查

支持 HTTP 和 gRPC 健康检查。

#### HTTP 健康检查

```go
import "github.com/ink-yht-code/sprout/sproutx/health"

checker := health.NewChecker()

// 注册依赖检查
checker.Register("database", func() error {
    return db.Ping()
})

checker.Register("redis", func() error {
    return redis.Ping(ctx).Err()
})

// 注册路由
engine.GET("/health", checker.HTTPHandler())
```

#### gRPC 健康检查

```go
import "google.golang.org/grpc/health"
import "google.golang.org/grpc/health/grpc_health_v1"

healthServer := health.NewServer()
grpc_health_v1.RegisterHealthServer(server.Server(), healthServer)

// 设置服务状态
healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_SERVING)
```

#### 健康检查响应

```json
{
    "status": "ok",
    "checks": {
        "database": "ok",
        "redis": "ok"
    }
}
```

---

### outbox - 事务消息

Outbox 模式，保证数据库事务和消息发送的一致性。

#### 创建 Outbox

```go
import "github.com/ink-yht-code/sprout/sproutx/outbox"

outboxStore := outbox.NewStore(db)
producer := outbox.NewProducer(outboxStore, kafkaProducer)
```

#### 使用 Outbox

```go
// 在事务中发送消息
err := txMgr.Do(ctx, func(ctx context.Context) error {
    // 保存业务数据
    if err := orderRepo.Save(ctx, order); err != nil {
        return err
    }
    
    // 发送消息（会在事务提交后发送）
    event := outbox.Event{
        AggregateID:   order.ID,
        AggregateType: "order",
        EventType:     "order.created",
        Payload:       order,
    }
    return producer.Send(ctx, event)
})
```

#### 后台发送

```go
// 启动后台发送器
sender := outbox.NewSender(outboxStore, kafkaProducer)
go sender.Run(ctx)
```

---

## 配置说明

### 完整配置结构

```go
type Config struct {
    Service ServiceConfig
    HTTP    HTTPConfig
    RPC     RPCConfig
    Log     LogConfig
    DB      DBConfig
    Redis   RedisConfig
}

type ServiceConfig struct {
    ID   int    // ServiceID
    Name string // 服务名称
}

type HTTPConfig struct {
    Enabled bool   // 是否启用
    Addr    string // 监听地址
    Mode    string // 运行模式：debug, release, test
}

type RPCConfig struct {
    Enabled bool   // 是否启用
    Addr    string // 监听地址
}

type LogConfig struct {
    Level    string     // 日志级别：debug, info, warn, error
    Encoding string     // 编码格式：json, console
    Output   string     // 输出位置：stdout, file
    File     FileConfig // 文件配置
}

type FileConfig struct {
    Path     string // 文件路径
    MaxSize  int    // 最大大小（MB）
    MaxAge   int    // 最大保留天数
    Compress bool   // 是否压缩
}

type DBConfig struct {
    DSN      string // 数据库连接字符串
    MaxOpen  int    // 最大打开连接数
    MaxIdle  int    // 最大空闲连接数
    LogLevel string // 日志级别：silent, error, warn, info
}

type RedisConfig struct {
    Addr         string // Redis 地址
    Password     string // 密码
    DB           int    // 数据库编号
    PoolSize     int    // 连接池大小
    MinIdleConns int    // 最小空闲连接数
}
```

### YAML 配置示例

```yaml
service:
  id: 101
  name: user

http:
  enabled: true
  addr: ":8080"
  mode: "release"

rpc:
  enabled: true
  addr: ":9090"

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
  pool_size: 100
  min_idle_conns: 10
```

---

## 许可证

Proprietary License

未经版权所有者书面授权，不得使用、复制、修改或分发本项目的任何部分。

---

Made with ❤️ by [ink-yht-code](https://github.com/ink-yht-code)
