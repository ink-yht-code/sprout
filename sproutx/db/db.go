package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 是数据库配置。
type Config struct {
	DSN      string
	MaxOpen  int
	MaxIdle  int
	LogLevel string
}

// New 创建并初始化 GORM DB。
//
// LogLevel 支持 silent/error/warn/info。
func New(cfg Config) (*gorm.DB, error) {
	var level logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	case "info":
		level = logger.Info
	default:
		level = logger.Info
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	}
	if cfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
