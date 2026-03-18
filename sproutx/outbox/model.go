package outbox

import (
	"time"

	"gorm.io/gorm"
)

// Status 表示 Outbox 记录状态。
type Status string

const (
	// StatusPending 表示待发送。
	StatusPending Status = "pending"
	// StatusSent 表示已发送。
	StatusSent Status = "sent"
	// StatusFailed 表示发送失败。
	StatusFailed Status = "failed"
)

// Outbox 表示事务消息记录。
type Outbox struct {
	ID           string    `gorm:"primaryKey;size:64"`
	Service      string    `gorm:"size:32;not null;index"`
	EventName    string    `gorm:"size:64;not null"`
	Payload      string    `gorm:"type:json;not null"`
	Status       Status    `gorm:"size:16;not null;default:pending;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	NextRetryAt  time.Time `gorm:"index"`
	RetryCount   int       `gorm:"default:0"`
	RequestID    string    `gorm:"size:64"`
	ErrorMessage string    `gorm:"size:512"`
}

// TableName 返回表名。
func (Outbox) TableName() string {
	return "outbox"
}

// Repository 定义 Outbox 仓储接口。
type Repository interface {
	Save(db *gorm.DB, records ...*Outbox) error
	ListPending(db *gorm.DB, limit int) ([]*Outbox, error)
	MarkSent(db *gorm.DB, ids ...string) error
	MarkFailed(db *gorm.DB, id string, errMsg string) error
}

// Repo 是 Repository 的默认实现。
type Repo struct{}

// NewRepo 创建 Repo。
func NewRepo() *Repo {
	return &Repo{}
}

// Save 保存一批 Outbox 记录。
func (r *Repo) Save(db *gorm.DB, records ...*Outbox) error {
	return db.Create(&records).Error
}

// ListPending 查询待发送且到达重试时间的记录。
func (r *Repo) ListPending(db *gorm.DB, limit int) ([]*Outbox, error) {
	var records []*Outbox
	err := db.Where("status = ? AND next_retry_at <= ?", StatusPending, time.Now()).
		Limit(limit).
		Find(&records).Error
	return records, err
}

// MarkSent 将指定记录标记为已发送。
func (r *Repo) MarkSent(db *gorm.DB, ids ...string) error {
	return db.Model(&Outbox{}).Where("id IN ?", ids).
		Updates(map[string]any{
			"status":  StatusSent,
			"sent_at": time.Now(),
		}).Error
}

// MarkFailed 将指定记录标记为发送失败，并递增重试次数。
func (r *Repo) MarkFailed(db *gorm.DB, id string, errMsg string) error {
	return db.Model(&Outbox{}).Where("id = ?", id).
		Updates(map[string]any{
			"status":        StatusFailed,
			"error_message": errMsg,
			"retry_count":   gorm.Expr("retry_count + 1"),
		}).Error
}
