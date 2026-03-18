package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ink-yht-code/sprout/sprout-registry/store"
)

// Handler 提供 ServiceID 注册相关的 HTTP 处理器。
type Handler struct {
	store store.Store
}

// NewHandler 创建 Handler 实例。
func NewHandler(store store.Store) *Handler {
	return &Handler{store: store}
}

// AllocateRequest 表示服务 ID 分配请求。
type AllocateRequest struct {
	Name string `json:"name" binding:"required"`
}

// AllocateResponse 表示服务 ID 分配响应。
type AllocateResponse struct {
	// Allocate 处理 POST /v1/services:allocate，为给定服务名分配或返回已有 ServiceID。
	ServiceID int    `json:"service_id"`
	Name      string `json:"name"`
}

func (h *Handler) Allocate(c *gin.Context) {
	var req AllocateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	svc, err := h.store.Allocate(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AllocateResponse{
		// Get 处理 GET /v1/services/:name，按名称查询服务信息。
		ServiceID: svc.ServiceID,
		Name:      svc.Name,
	})
}

func (h *Handler) Get(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}

	svc, err := h.store.Get(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	c.JSON(http.StatusOK, AllocateResponse{
		ServiceID: svc.ServiceID,
		Name:      svc.Name,
	})
}

func (h *Handler) List(c *gin.Context) {
	services, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]AllocateResponse, len(services))
	for i, svc := range services {
		result[i] = AllocateResponse{
			ServiceID: svc.ServiceID,
			Name:      svc.Name,
		}
	}

	c.JSON(http.StatusOK, gin.H{"services": result})
}

// RegisterRoutes 将路由注册到 Gin 引擎。
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.POST("/services:allocate", h.Allocate)
		v1.GET("/services/:name", h.Get)
		v1.GET("/services", h.List)
	}
}
