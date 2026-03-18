package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var _ Manager = (*manager)(nil)

type manager struct {
	opts Options
}

// NewManager 创建一个 JWT 管理器。
//
// 通过 Options 指定签名密钥、过期时间、签名算法以及 Issuer。
func NewManager(opts Options) Manager {
	return &manager{
		opts: opts,
	}
}

func (m *manager) GenerateToken(claims Claims) (string, error) {
	return m.generateToken(claims, m.opts.AccessExpire)
}

func (m *manager) GenerateTokenPair(claims Claims) (*TokenPair, error) {
	accessToken, err := m.generateToken(claims, m.opts.AccessExpire)
	if err != nil {
		return nil, fmt.Errorf("生成 Access Token 失败: %w", err)
	}

	refreshToken, err := m.generateToken(claims, m.opts.RefreshExpire)
	if err != nil {
		return nil, fmt.Errorf("生成 Refresh Token 失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *manager) generateToken(claims Claims, expire time.Duration) (string, error) {
	now := time.Now()

	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    m.opts.Issuer,
		ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ID:        uuid.New().String(),
	}

	token := jwt.NewWithClaims(m.opts.Method, claims)

	return token.SignedString([]byte(m.opts.SignKey))
}

func (m *manager) VerifyToken(tokenString string) (*Claims, error) {
	return m.verifyToken(tokenString)
}

func (m *manager) VerifyRefreshToken(tokenString string) (*Claims, error) {
	return m.verifyToken(tokenString)
}

func (m *manager) verifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != m.opts.Method {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(m.opts.SignKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析 Token 失败: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("无效的 Token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("无效的 Claims 类型")
	}

	return claims, nil
}
