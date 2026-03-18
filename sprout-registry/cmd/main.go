package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ink-yht-code/sprout/sprout-registry/api"
	"github.com/ink-yht-code/sprout/sprout-registry/store"
)

var (
	port     = flag.Int("port", 18080, "HTTP port")
	dataFile = flag.String("data", "registry.db", "SQLite database file")
	token    = flag.String("token", "", "Authentication token")
)

func main() {
	flag.Parse()

	s, err := store.NewSQLiteStore(*dataFile)
	if err != nil {
		fmt.Printf("Failed to init store: %v\n", err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	if *token != "" {
		r.Use(func(c *gin.Context) {
			auth := c.GetHeader("Authorization")
			if auth != "Bearer "+*token {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				c.Abort()
				return
			}
			c.Next()
		})
	}

	h := api.NewHandler(s)
	h.RegisterRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "sprout",
			"description": "基于 Gin 构建的 Go 微服务框架，提供路由、认证、校验等开箱即用的功能",
			"version":     "1.0.0",
		})
	})

	addr := fmt.Sprintf(":%d", *port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		fmt.Printf("Registry server listening on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exited")
}
