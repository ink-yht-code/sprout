package tx

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey struct{}

// Manager 是事务管理器。
type Manager struct {
	db *gorm.DB
}

// NewManager 创建事务管理器。
func NewManager(db *gorm.DB) *Manager {
	return &Manager{db: db}
}

// Do 在一个事务内执行 fn。
//
// 事务 DB 会被写入 context，供 DAO/Repository 通过 FromContext 获取。
func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxKey{}, tx)
		return fn(ctx)
	})
}

// FromContext 从 context 读取事务 DB，若不存在则返回 defaultDB。
func FromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if ctx == nil {
		return defaultDB
	}
	if tx, ok := ctx.Value(ctxKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}

// GetDB 是 FromContext 的别名。
func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	return FromContext(ctx, defaultDB)
}
