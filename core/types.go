package core

import (
	"github.com/gin-gonic/gin"
)

// Handler 定义了一个 HTTP Handler 的路由注册约定。
//
// PublicRoutes 用于注册无需鉴权的公开路由；PrivateRoutes 用于注册需要鉴权的私有路由。
type Handler interface {
	PrivateRoutes(server *gin.Engine)
	PublicRoutes(server *gin.Engine)
}

// Result 是统一的 HTTP JSON 响应结构。
//
// Code 表示业务码；Msg 表示提示信息；Data 表示响应数据。
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// PageData 表示分页返回数据。
//
// List 为数据列表；Total 为总数；Page/Size 为当前页码和每页大小。
type PageData[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

// PageRequest 表示分页请求参数。
type PageRequest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

// Validate 校验并修正分页参数。
//
// Page 最小为 1；Size 默认 10，最大 100。
func (p *PageRequest) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}
}

// Offset 计算分页的偏移量（用于数据库查询）。
func (p *PageRequest) Offset() int {
	return (p.Page - 1) * p.Size
}
