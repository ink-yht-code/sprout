package store

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Service 表示已注册服务的数据库模型。
type Service struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;size:64;not null"`
	ServiceID int       `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 返回 Service 对应的数据库表名。
func (Service) TableName() string {
	return "services"
}

// Store 定义服务注册存储的接口。
type Store interface {
	Allocate(name string) (*Service, error)
	Get(name string) (*Service, error)
	List() ([]Service, error)
}

// SQLiteStore 基于 SQLite 的 Store 实现。
type SQLiteStore struct {
	db *gorm.DB
}

// NewSQLiteStore 创建 SQLiteStore 实例并执行自动迁移。
func NewSQLiteStore(dsn string) (*SQLiteStore, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Service{}); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Allocate(name string) (*Service, error) {
	var svc Service
	if err := s.db.Where("name = ?", name).First(&svc).Error; err == nil {
		return &svc, nil
	}

	var maxID int
	s.db.Model(&Service{}).Select("COALESCE(MAX(service_id), 100)").Scan(&maxID)
	nextID := maxID + 1
	if nextID < 101 {
		nextID = 101
	}

	svc = Service{
		Name:      name,
		ServiceID: nextID,
	}
	if err := s.db.Create(&svc).Error; err != nil {
		return nil, err
	}

	return &svc, nil
}

func (s *SQLiteStore) Get(name string) (*Service, error) {
	var svc Service
	if err := s.db.Where("name = ?", name).First(&svc).Error; err != nil {
		return nil, err
	}
	return &svc, nil
}

// List 返回所有已注册服务。
func (s *SQLiteStore) List() ([]Service, error) {
	var services []Service
	if err := s.db.Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}
