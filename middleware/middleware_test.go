package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestTokenBucketLimiter(t *testing.T) {
	limiter := NewTokenBucketLimiter(5, time.Second)

	for i := 0; i < 5; i++ {
		if !limiter.Allow("test-key") {
			t.Errorf("Expected allow on request %d", i+1)
		}
	}

	if limiter.Allow("test-key") {
		t.Error("Expected deny on 6th request")
	}

	time.Sleep(1100 * time.Millisecond)

	if !limiter.Allow("test-key") {
		t.Error("Expected allow after refill")
	}
}

func TestRateLimitByIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewTokenBucketLimiter(3, time.Second)
	router := gin.New()
	router.Use(RateLimitByIP(limiter, []string{"/health"}))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i+1, w.Code)
		}
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}
}

func TestRateLimitExcludedPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewTokenBucketLimiter(1, time.Second)
	router := gin.New()
	router.Use(RateLimitByIP(limiter, []string{"/health", "/ping"}))
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i+1, w.Code)
		}
	}
}

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		handler        gin.HandlerFunc
		expectedStatus int
	}{
		{
			name: "normal handler",
			handler: func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "panic handler",
			handler: func(c *gin.Context) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(Recovery(DefaultRecoveryConfig()))
			router.GET("/test", tt.handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if len(config.CORS.AllowOrigins) == 0 {
		t.Error("Expected CORS config")
	}

	if len(config.CSRF.TokenLookup) == 0 {
		t.Error("Expected CSRF config")
	}

	if config.Recovery.StackAll != false {
		t.Error("Expected StackAll to be false")
	}
}

func TestNewDefaultMiddleware(t *testing.T) {
	middlewares := NewDefaultMiddleware()

	if len(middlewares) == 0 {
		t.Error("Expected at least one middleware")
	}
}

func TestNewRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewRateLimitMiddleware(10)
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i+1, w.Code)
		}
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}
}
