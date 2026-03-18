package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ink-yht-code/sprout/jwt"
	"github.com/ink-yht-code/sprout/session"
)

// JWTAuthConfig 是 JWT 鉴权中间件配置。
type JWTAuthConfig struct {
	Manager jwt.Manager
	Header  string
	Prefix  string
}

// DefaultJWTAuthConfig 返回默认 JWT 鉴权配置。
func DefaultJWTAuthConfig(manager jwt.Manager) JWTAuthConfig {
	return JWTAuthConfig{
		Manager: manager,
		Header:  "Authorization",
		Prefix:  "Bearer ",
	}
}

// JWTAuth 创建 JWT 鉴权中间件。
//
// 校验成功后会往 gin.Context 写入：
// - user_id
// - session.CtxSessionIDKey（SSID）
func JWTAuth(cfg JWTAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.GetHeader(cfg.Header)
		if val == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(val, cfg.Prefix) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(val, cfg.Prefix))
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := cfg.Manager.VerifyToken(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set(session.CtxSessionIDKey, claims.SSID)
		c.Next()
	}
}
