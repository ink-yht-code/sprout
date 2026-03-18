package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisProvider 是基于 Redis 的 Session Provider 实现。
//
// Session 会被序列化为 JSON 存入 Redis，并以 TTL 控制过期。
type RedisProvider struct {
	client *redis.Client
	prefix string
}

// NewRedisProvider 创建一个 Redis Provider。
//
// prefix 用于区分不同应用的 key 命名空间。
func NewRedisProvider(client *redis.Client, prefix string) Provider {
	return &RedisProvider{
		client: client,
		prefix: prefix,
	}
}

// Create 创建并保存一个 Session。
func (p *RedisProvider) Create(ctx context.Context, sessionID string, data map[string]interface{}, expire time.Duration) error {
	sess := &Session{
		ID:        sessionID,
		Data:      data,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expire),
	}

	dataBytes, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("序列化 session 失败: %w", err)
	}

	key := p.key(sessionID)
	return p.client.Set(ctx, key, dataBytes, expire).Err()
}

// Get 获取 Session。
func (p *RedisProvider) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := p.key(sessionID)
	data, err := p.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("获取 session 失败: %w", err)
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("反序列化 session 失败: %w", err)
	}

	if time.Now().After(sess.ExpiresAt) {
		p.client.Del(ctx, key)
		return nil, ErrSessionExpired
	}

	return &sess, nil
}

// Update 更新 Session 数据（保持原过期时间）。
func (p *RedisProvider) Update(ctx context.Context, sessionID string, data map[string]interface{}) error {
	sess, err := p.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	sess.Data = data
	dataBytes, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("序列化 session 失败: %w", err)
	}

	key := p.key(sessionID)
	expire := time.Until(sess.ExpiresAt)
	if expire <= 0 {
		return ErrSessionExpired
	}

	return p.client.Set(ctx, key, dataBytes, expire).Err()
}

// Delete 删除 Session。
func (p *RedisProvider) Delete(ctx context.Context, sessionID string) error {
	key := p.key(sessionID)
	return p.client.Del(ctx, key).Err()
}

// Refresh 刷新 Session 的过期时间。
func (p *RedisProvider) Refresh(ctx context.Context, sessionID string, expire time.Duration) error {
	sess, err := p.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	sess.ExpiresAt = time.Now().Add(expire)
	dataBytes, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("序列化 session 失败: %w", err)
	}

	key := p.key(sessionID)
	return p.client.Set(ctx, key, dataBytes, expire).Err()
}

// Close 关闭 Redis 客户端连接。
func (p *RedisProvider) Close() error {
	return p.client.Close()
}

func (p *RedisProvider) key(sessionID string) string {
	return fmt.Sprintf("%s:%s", p.prefix, sessionID)
}
