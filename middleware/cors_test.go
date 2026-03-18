package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		config         CORSConfig
		origin         string
		method         string
		expectedStatus int
		expectHeaders  bool
	}{
		{
			name: "allow all origins",
			config: CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"GET", "POST"},
				AllowHeaders: []string{"Content-Type"},
			},
			origin:         "http://example.com",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectHeaders:  true,
		},
		{
			name: "specific origin",
			config: CORSConfig{
				AllowOrigins: []string{"http://example.com"},
				AllowMethods: []string{"GET"},
				AllowHeaders: []string{"Content-Type"},
			},
			origin:         "http://example.com",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectHeaders:  true,
		},
		{
			name: "disallowed origin",
			config: CORSConfig{
				AllowOrigins: []string{"http://example.com"},
				AllowMethods: []string{"GET"},
				AllowHeaders: []string{"Content-Type"},
			},
			origin:         "http://other.com",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectHeaders:  false,
		},
		{
			name: "options request",
			config: CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"GET", "POST"},
				AllowHeaders: []string{"Content-Type"},
			},
			origin:         "http://example.com",
			method:         "OPTIONS",
			expectedStatus: http.StatusNoContent,
			expectHeaders:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(CORS(tt.config))
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			})

			req := httptest.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectHeaders {
				if w.Header().Get("Access-Control-Allow-Origin") == "" {
					t.Error("Expected Access-Control-Allow-Origin header")
				}
				if w.Header().Get("Access-Control-Allow-Methods") == "" {
					t.Error("Expected Access-Control-Allow-Methods header")
				}
				if w.Header().Get("Access-Control-Allow-Headers") == "" {
					t.Error("Expected Access-Control-Allow-Headers header")
				}
			}
		})
	}
}

func TestDefaultCORSConfig(t *testing.T) {
	config := DefaultCORSConfig()

	if len(config.AllowOrigins) == 0 {
		t.Error("Expected non-empty AllowOrigins")
	}

	if len(config.AllowMethods) == 0 {
		t.Error("Expected non-empty AllowMethods")
	}

	if len(config.AllowHeaders) == 0 {
		t.Error("Expected non-empty AllowHeaders")
	}

	if config.MaxAge != 86400 {
		t.Errorf("Expected MaxAge 86400, got %d", config.MaxAge)
	}
}
