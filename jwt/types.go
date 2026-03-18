package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 是 Sprout 使用的 JWT Claims。
//
// UserId 为用户标识；SSID 为会话标识；Data 为自定义键值对。
type Claims struct {
	UserId string            `json:"user_id"`
	SSID   string            `json:"ssid"`
	Data   map[string]string `json:"data"`
	jwt.RegisteredClaims
}

// Options 是 JWT 管理器的配置项。
type Options struct {
	SignKey       string
	AccessExpire  time.Duration
	RefreshExpire time.Duration
	Method        jwt.SigningMethod
	Issuer        string
}

// NewOptions 创建默认的 JWT 配置。
//
// 默认签名方法为 HS256，Issuer 为 sprout。
func NewOptions(signKey string, accessExpire, refreshExpire time.Duration) Options {
	return Options{
		SignKey:       signKey,
		AccessExpire:  accessExpire,
		RefreshExpire: refreshExpire,
		Method:        jwt.SigningMethodHS256,
		Issuer:        "sprout",
	}
}

// TokenPair 表示一对访问令牌与刷新令牌。
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Manager 定义 JWT Token 的生成与校验接口。
type Manager interface {
	GenerateToken(claims Claims) (string, error)
	GenerateTokenPair(claims Claims) (*TokenPair, error)
	VerifyToken(token string) (*Claims, error)
	VerifyRefreshToken(token string) (*Claims, error)
}
