package scaffoldtmpl

// WiringTmpl wiring 模板
var WiringTmpl = `package wiring

import (
	"context"

	"github.com/ink-yht-code/sprout/jwt"
	"github.com/ink-yht-code/sprout/sproutx/httpx"
	"github.com/ink-yht-code/sprout/sproutx/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"{{.Name}}/internal/config"
	"{{.Name}}/internal/repository"
	"{{.Name}}/internal/repository/cache"
	"{{.Name}}/internal/repository/dao"
	"{{.Name}}/internal/service"
	"{{.Name}}/internal/web"
)

// LoadConfig 加载配置
func LoadConfig(path string) (*config.Config, error) {
	return config.Load(path)
}

// App 应用
type App struct {
	cfg  *config.Config
	http *httpx.Server
	db   *gorm.DB
}

// BuildApp 构建应用
func BuildApp(cfg *config.Config) (*App, error) {
	// 初始化日志
	if err := log.Init(log.Config{
		Level:    cfg.Log.Level,
		Encoding: cfg.Log.Encoding,
		Output:   cfg.Log.Output,
	}); err != nil {
		return nil, err
	}

	// 初始化 JWT Manager
	jwtManager := jwt.NewManager(jwt.Options{
		SignKey:        cfg.JWT.Secret,
		AccessExpire:  cfg.JWT.AccessExpire,
		RefreshExpire: cfg.JWT.RefreshExpire,
		Issuer:        cfg.JWT.Issuer,
	})

	// 初始化数据库 (可选)
	var db *gorm.DB
	if cfg.DB.DSN != "" {
		var err error
		db, err = gorm.Open(mysql.Open(cfg.DB.DSN), &gorm.Config{})
		if err != nil {
			log.Warn("Failed to connect database, using memory DAO", zap.Error(err))
		}
	}

	// 初始化 Redis (可选)
	var rdb redis.Cmdable
	if cfg.Redis.Addr != "" {
		rdb = redis.NewClient(&redis.Options{
			Addr: cfg.Redis.Addr,
		})
	}

	// 创建 DAO 和 Repository
	var d dao.{{.NameUpper}}DAO = dao.New{{.NameUpper}}DAO(db)
	var c cache.{{.NameUpper}}Cache = cache.New{{.NameUpper}}Cache(rdb)
	repo := repository.New{{.NameUpper}}Repository(d, c)
	svc := service.New{{.NameUpper}}Service(repo, jwtManager)

	// 创建 Handler
	handler := web.NewHandler(svc)

	// 创建 HTTP server
	var httpServer *httpx.Server
	if cfg.HTTP.Enabled {
		httpServer = httpx.NewServer(httpx.Config{
			Enabled: cfg.HTTP.Enabled,
			Addr:    cfg.HTTP.Addr,
		})
		
		// 注册公开路由
		handler.PublicRoutes(httpServer.Engine)
		
		// 注册私有路由 (需要认证)
		// 可以使用 JWT 中间件保护私有路由
		// authGroup := httpServer.Engine.Group("/")
		// authGroup.Use(jwtMiddleware)
		// handler.PrivateRoutes(httpServer.Engine)
		
		// 框架介绍
		httpServer.Engine.GET("/", handler.Health)
	}

	return &App{cfg: cfg, http: httpServer, db: db}, nil
}

// Run 启动应用
func (a *App) Run() error {
	if a.http != nil {
		log.Info("HTTP server starting", zap.String("addr", a.cfg.HTTP.Addr))
		return a.http.Run()
	}
	return nil
}

// Shutdown 关闭应用
func (a *App) Shutdown(ctx context.Context) error {
	if a.http != nil {
		return a.http.Shutdown(ctx)
	}
	return nil
}
`
