package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ink-yht-code/sprout/sprout-registry/api"
	"github.com/ink-yht-code/sprout/sprout-registry/store"
	"github.com/spf13/cobra"
)

// NewCommand 创建 registry 子命令，用于启动 ServiceID 注册服务。
func NewCommand() *cobra.Command {
	var (
		port     int
		dataFile string
		token    string
	)

	cmd := &cobra.Command{
		Use:   "registry",
		Short: "ServiceID 注册服务",
		Long:  "集中式 ServiceID 分配服务，确保微服务 ID 全局唯一",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := store.NewSQLiteStore(dataFile)
			if err != nil {
				return fmt.Errorf("failed to init store: %w", err)
			}

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.Use(gin.Recovery())

			if token != "" {
				r.Use(func(c *gin.Context) {
					auth := c.GetHeader("Authorization")
					if auth != "Bearer "+token {
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

			addr := fmt.Sprintf(":%d", port)
			srv := &http.Server{Addr: addr, Handler: r}

			go func() {
				cmd.Printf("Registry server listening on %s\n", addr)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					cmd.PrintErrf("Failed to start server: %v\n", err)
					os.Exit(1)
				}
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			cmd.Println("Shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				cmd.PrintErrf("Server forced to shutdown: %v\n", err)
			}

			cmd.Println("Server exited")
			return nil
		},
	}

	cmd.Flags().IntVar(&port, "port", 18080, "HTTP port")
	cmd.Flags().StringVar(&dataFile, "data", "registry.db", "SQLite database file")
	cmd.Flags().StringVar(&token, "token", "", "Authentication token")

	return cmd
}
