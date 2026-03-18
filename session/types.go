package session

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ink-yht-code/sprout/jwt"
)

// ErrSessionNotFound 表示未找到 Session。
var ErrSessionNotFound = errors.New("session not found")

// ErrSessionExpired 表示 Session 已过期。
var ErrSessionExpired = errors.New("session expired")

// ErrProviderNotInitialized 表示未设置默认 Session Provider。
var ErrProviderNotInitialized = errors.New("session provider 未初始化")

// CtxSessionIDKey 是在 Context 中保存 SessionID 的键。
const CtxSessionIDKey = "sprout:session_id"

// Session 表示一次会话的数据。
//
// ID 为会话 ID；Data 为会话数据；CreatedAt/ExpiresAt 为创建与过期时间；JwtClaims 为绑定在该会话上的 JWT Claims。
type Session struct {
	ID        string
	Data      map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
	JwtClaims *jwt.Claims
}

// Claims 返回会话绑定的 JWT Claims。
func (s Session) Claims() *jwt.Claims {
	return s.JwtClaims
}

// Get 获取会话数据。
func (s *Session) Get(key string) (interface{}, bool) {
	if s.Data == nil {
		return nil, false
	}
	val, ok := s.Data[key]
	return val, ok
}

// Set 设置会话数据。
func (s *Session) Set(key string, val interface{}) {
	if s.Data == nil {
		s.Data = make(map[string]interface{})
	}
	s.Data[key] = val
}

// Del 删除会话数据。
func (s *Session) Del(key string) {
	if s.Data != nil {
		delete(s.Data, key)
	}
}

// Provider 定义 Session 存储接口。
type Provider interface {
	Create(ctx context.Context, sessionID string, data map[string]interface{}, expire time.Duration) error
	Get(ctx context.Context, sessionID string) (*Session, error)
	Update(ctx context.Context, sessionID string, data map[string]interface{}) error
	Delete(ctx context.Context, sessionID string) error
	Refresh(ctx context.Context, sessionID string, expire time.Duration) error
	Close() error
}

// TokenCarrier 描述 Token 的注入、提取与清理方式。
type TokenCarrier interface {
	Inject(ctx context.Context, token string)
	Extract(ctx context.Context) string
	Clear(ctx context.Context)
}

// HeaderTokenCarrier 使用 Authorization Header 传递 Token。
type HeaderTokenCarrier struct{}

// Inject 将 Token 写入响应头（Authorization: Bearer <token>）。
func (h *HeaderTokenCarrier) Inject(ctx context.Context, token string) {
	if gc, ok := ctx.(*ginContext); ok {
		gc.Writer.Header().Set("Authorization", "Bearer "+token)
	}
}

// Extract 从请求头或响应头读取 Authorization。
func (h *HeaderTokenCarrier) Extract(ctx context.Context) string {
	if gc, ok := ctx.(*ginContext); ok {
		auth := gc.Writer.Header().Get("Authorization")
		if auth != "" {
			return auth
		}
		auth = gc.Request.Header.Get("Authorization")
		if auth != "" {
			return auth
		}
	}
	return ""
}

// Clear 清理响应头中的 Authorization。
func (h *HeaderTokenCarrier) Clear(ctx context.Context) {
	if gc, ok := ctx.(*ginContext); ok {
		gc.Writer.Header().Del("Authorization")
	}
}

// CookieTokenCarrier 使用 Cookie 传递 Token。
type CookieTokenCarrier struct {
	Name     string
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite int
}

// Inject 将 Token 写入 Cookie。
func (c *CookieTokenCarrier) Inject(ctx context.Context, token string) {
	if gc, ok := ctx.(*ginContext); ok {
		gc.SetCookie(c.Name, token, 0, c.Path, c.Domain, c.Secure, c.HttpOnly)
	}
}

// Extract 从 Cookie 中读取 Token。
func (c *CookieTokenCarrier) Extract(ctx context.Context) string {
	if gc, ok := ctx.(*ginContext); ok {
		cookie, err := gc.Cookie(c.Name)
		if err == nil {
			return cookie
		}
	}
	return ""
}

// Clear 清理 Cookie 中的 Token。
func (c *CookieTokenCarrier) Clear(ctx context.Context) {
	if gc, ok := ctx.(*ginContext); ok {
		gc.SetCookie(c.Name, "", -1, c.Path, c.Domain, c.Secure, c.HttpOnly)
	}
}

type ginContext struct {
	*gin.Context
}

var defaultProvider atomic.Value

// SetDefaultProvider 设置默认的 Session Provider。
func SetDefaultProvider(provider Provider) {
	defaultProvider.Store(provider)
}

func getDefaultProvider() (Provider, error) {
	p := defaultProvider.Load()
	if p == nil {
		return nil, ErrProviderNotInitialized
	}
	return p.(Provider), nil
}

// Get 获取当前请求上下文中的 Session。
//
// SessionID 需要先写入 Context（key 为 CtxSessionIDKey）。
func Get(ctx context.Context) (Session, error) {
	p, err := getDefaultProvider()
	if err != nil {
		return Session{}, err
	}

	var sid string
	if v, ok := ctx.(interface{ Get(string) (any, bool) }); ok {
		if val, ok2 := v.Get(CtxSessionIDKey); ok2 {
			sid, _ = val.(string)
		}
	}
	if sid == "" {
		return Session{}, ErrSessionNotFound
	}

	sess, err := p.Get(ctx, sid)
	if err != nil {
		return Session{}, err
	}
	if sess == nil {
		return Session{}, ErrSessionNotFound
	}
	return *sess, nil
}

// NewSession 创建一个新的 Session，并将 SessionID 写入 Context。
//
// 默认过期时间为 24 小时。
func NewSession(ctx context.Context, userId string, jwtData map[string]string, sessData map[string]any) (Session, error) {
	p, err := getDefaultProvider()
	if err != nil {
		return Session{}, err
	}

	sid := uuid.NewString()
	claims := &jwt.Claims{
		UserId: userId,
		SSID:   sid,
		Data:   jwtData,
	}

	data := make(map[string]interface{}, len(sessData))
	for k, v := range sessData {
		data[k] = v
	}

	expire := 24 * time.Hour
	if err := p.Create(ctx, sid, data, expire); err != nil {
		return Session{}, err
	}

	sess := Session{
		ID:        sid,
		Data:      data,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expire),
		JwtClaims: claims,
	}

	if v, ok := ctx.(interface{ Set(string, any) }); ok {
		v.Set(CtxSessionIDKey, sid)
	}

	return sess, nil
}
