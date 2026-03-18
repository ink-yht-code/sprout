package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CSRFConfig 是 CSRF 防护中间件配置。
type CSRFConfig struct {
	TokenLength    int
	TokenLookup    string
	CookieName     string
	CookieDomain   string
	CookiePath     string
	CookieSecure   bool
	CookieHTTPOnly bool
	HeaderName     string
	FormName       string
	TrustedOrigins []string
}

// DefaultCSRFConfig 返回默认 CSRF 配置。
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLength:    32,
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "_csrf",
		CookieDomain:   "",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: false,
		HeaderName:     "X-CSRF-Token",
		FormName:       "_csrf",
		TrustedOrigins: []string{},
	}
}

// CSRF 创建 CSRF 防护中间件。
//
// 对 GET/HEAD/OPTIONS 请求不做校验；对其他请求校验 Origin/Referer 以及 cookie token 与客户端 token 是否一致。
func CSRF(config CSRFConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		origin := c.Request.Header.Get("Origin")
		referer := c.Request.Header.Get("Referer")

		if !isTrustedOrigin(origin, referer, config.TrustedOrigins) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF origin not trusted"})
			return
		}

		cookieToken, err := c.Cookie(config.CookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF cookie not found"})
			return
		}

		var clientToken string
		if strings.HasPrefix(config.TokenLookup, "header:") {
			headerName := strings.TrimPrefix(config.TokenLookup, "header:")
			clientToken = c.GetHeader(headerName)
		} else if strings.HasPrefix(config.TokenLookup, "form:") {
			formName := strings.TrimPrefix(config.TokenLookup, "form:")
			clientToken = c.PostForm(formName)
		}

		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token not provided"})
			return
		}

		if subtle.ConstantTimeCompare([]byte(cookieToken), []byte(clientToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			return
		}

		c.Next()
	}
}

func isTrustedOrigin(origin, referer string, trustedOrigins []string) bool {
	if len(trustedOrigins) == 0 {
		return true
	}

	for _, trusted := range trustedOrigins {
		if origin == trusted || strings.HasPrefix(referer, trusted) {
			return true
		}
	}

	return false
}
