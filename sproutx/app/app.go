package app

import (
	"context"
	"fmt"

	"github.com/ink-yht-code/sprout/sproutx/db"
	"github.com/ink-yht-code/sprout/sproutx/httpx"
	"github.com/ink-yht-code/sprout/sproutx/log"
	"github.com/ink-yht-code/sprout/sproutx/redis"
	"github.com/ink-yht-code/sprout/sproutx/tx"
	redislib "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Config 是应用启动器配置。
type Config struct {
	Service ServiceConfig
	HTTP    httpx.Config
	Log     log.Config
	DB      db.Config
	Redis   redis.Config
}

// ServiceConfig 是服务基础信息配置。
type ServiceConfig struct {
	ID   int
	Name string
}

// App 是运行时应用实例，聚合 DB/Redis/事务/HTTP 等组件。
type App struct {
	Config    *Config
	DB        *gorm.DB
	Redis     *redislib.Client
	TxManager *tx.Manager
	HTTP      *httpx.Server
}

// New 创建并初始化 App。
//
// 会按配置初始化日志、数据库、Redis、HTTP Server 等组件。
func New(cfg *Config) (*App, error) {
	if err := log.Init(cfg.Log); err != nil {
		return nil, fmt.Errorf("init log: %w", err)
	}

	app := &App{Config: cfg}

	if cfg.DB.DSN != "" {
		database, err := db.New(cfg.DB)
		if err != nil {
			return nil, fmt.Errorf("init db: %w", err)
		}
		app.DB = database
		app.TxManager = tx.NewManager(database)
	}

	if cfg.Redis.Addr != "" {
		app.Redis = redislib.NewClient(&redislib.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	if cfg.HTTP.Enabled {
		app.HTTP = httpx.NewServer(cfg.HTTP)
	}

	return app, nil
}

// Run 启动应用（当前仅启动 HTTP Server）。
func (a *App) Run() error {
	if a.HTTP != nil {
		log.Info("HTTP server starting", zap.String("addr", a.Config.HTTP.Addr))
		if err := a.HTTP.Run(); err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	}
	return nil
}

// Shutdown 优雅关闭应用资源。
func (a *App) Shutdown(ctx context.Context) error {
	if a.HTTP != nil {
		if err := a.HTTP.Shutdown(ctx); err != nil {
			return err
		}
	}

	if a.DB != nil {
		sqlDB, _ := a.DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	if a.Redis != nil {
		a.Redis.Close()
	}

	return log.Sync()
}
