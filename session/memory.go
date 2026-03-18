package session

import (
	"context"
	"sync"
	"time"
)

// MemoryProvider 是基于内存的 Session Provider 实现。
//
// 适用于本地开发或单实例服务；进程退出后会丢失所有会话。
type MemoryProvider struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	stopChan chan struct{}
}

// NewMemoryProvider 创建一个内存 Provider，并启动后台清理过期 Session 的任务。
func NewMemoryProvider() Provider {
	p := &MemoryProvider{
		sessions: make(map[string]*Session),
		stopChan: make(chan struct{}),
	}
	go p.cleanup()
	return p
}

// Create 创建并保存一个 Session。
func (p *MemoryProvider) Create(ctx context.Context, sessionID string, data map[string]interface{}, expire time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.sessions[sessionID] = &Session{
		ID:        sessionID,
		Data:      data,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expire),
	}

	return nil
}

// Get 获取 Session。
func (p *MemoryProvider) Get(ctx context.Context, sessionID string) (*Session, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	sess, ok := p.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return sess, nil
}

// Update 更新 Session 数据。
func (p *MemoryProvider) Update(ctx context.Context, sessionID string, data map[string]interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	sess, ok := p.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	sess.Data = data
	return nil
}

// Delete 删除 Session。
func (p *MemoryProvider) Delete(ctx context.Context, sessionID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.sessions, sessionID)
	return nil
}

// Refresh 刷新 Session 的过期时间。
func (p *MemoryProvider) Refresh(ctx context.Context, sessionID string, expire time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	sess, ok := p.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	sess.ExpiresAt = time.Now().Add(expire)
	return nil
}

func (p *MemoryProvider) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.mu.Lock()
			now := time.Now()
			for id, sess := range p.sessions {
				if now.After(sess.ExpiresAt) {
					delete(p.sessions, id)
				}
			}
			p.mu.Unlock()
		case <-p.stopChan:
			return
		}
	}
}

// Close 关闭 Provider 并停止后台清理任务。
func (p *MemoryProvider) Close() error {
	close(p.stopChan)
	return nil
}
