package context

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Context 是对 gin.Context 的轻量包装，提供链式读取参数等便捷方法。
type Context struct {
	*gin.Context
}

// Value 是一个通用的参数读取结果，包含原始字符串值与可能的错误。
//
// 你可以使用 String/Int/Int64/Bool 系列方法做类型转换，也可以使用 *Or 方法提供默认值。
type Value struct {
	val string
	err error
}

// String 返回字符串值和错误。
func (v Value) String() (string, error) {
	return v.val, v.err
}

// StringOr 在出错时返回默认值。
func (v Value) StringOr(defaultVal string) string {
	if v.err != nil {
		return defaultVal
	}
	return v.val
}

// Int 将值转换为 int。
func (v Value) Int() (int, error) {
	if v.err != nil {
		return 0, v.err
	}
	return strconv.Atoi(v.val)
}

// IntOr 将值转换为 int，失败时返回默认值。
func (v Value) IntOr(defaultVal int) int {
	if v.err != nil {
		return defaultVal
	}
	val, err := strconv.Atoi(v.val)
	if err != nil {
		return defaultVal
	}
	return val
}

// Int64 将值转换为 int64。
func (v Value) Int64() (int64, error) {
	if v.err != nil {
		return 0, v.err
	}
	return strconv.ParseInt(v.val, 10, 64)
}

// Int64Or 将值转换为 int64，失败时返回默认值。
func (v Value) Int64Or(defaultVal int64) int64 {
	if v.err != nil {
		return defaultVal
	}
	val, err := strconv.ParseInt(v.val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

// Bool 将值转换为 bool。
func (v Value) Bool() (bool, error) {
	if v.err != nil {
		return false, v.err
	}
	return strconv.ParseBool(v.val)
}

// BoolOr 将值转换为 bool，失败时返回默认值。
func (v Value) BoolOr(defaultVal bool) bool {
	if v.err != nil {
		return defaultVal
	}
	val, err := strconv.ParseBool(v.val)
	if err != nil {
		return defaultVal
	}
	return val
}

// Param 获取路径参数（如 /users/:id 中的 id）。
func (c *Context) Param(key string) Value {
	return Value{
		val: c.Context.Param(key),
	}
}

// Query 获取查询参数（如 ?page=1 中的 page）。
func (c *Context) Query(key string) Value {
	return Value{
		val: c.Context.Query(key),
	}
}

// Cookie 获取 Cookie 值。
func (c *Context) Cookie(key string) Value {
	val, err := c.Context.Cookie(key)
	return Value{
		val: val,
		err: err,
	}
}

// Header 获取请求头的值。
func (c *Context) Header(key string) Value {
	return Value{
		val: c.Context.GetHeader(key),
	}
}

// UserId 从 Context 中读取 user_id（通常由鉴权中间件注入）。
func (c *Context) UserId() string {
	val, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	userId, ok := val.(string)
	if !ok {
		return ""
	}
	return userId
}

// SetUserId 将 user_id 写入 Context（通常在鉴权通过后调用）。
func (c *Context) SetUserId(userId string) {
	c.Set("user_id", userId)
}

// EventStream 创建 SSE（Server-Sent Events）事件通道，并将必要的 Header 写入响应。
//
// 返回的 channel 用于写入事件数据；当客户端断开时后台 goroutine 会自动退出。
func (c *Context) EventStream() chan []byte {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	eventCh := make(chan []byte, 10)

	go func() {
		for {
			select {
			case eventData, ok := <-eventCh:
				if !ok {
					return
				}
				if len(eventData) > 0 {
					c.sendEvent(eventData)
				}
			case <-c.Request.Context().Done():
				return
			}
		}
	}()

	return eventCh
}

func (c *Context) sendEvent(data []byte) {
	_, _ = c.Writer.Write(data)
	c.Writer.Flush()
}

// JSON 返回 JSON 响应（保持与 gin.Context.JSON 一致）。
func (c *Context) JSON(code int, obj any) {
	c.Context.JSON(code, obj)
}

// Success 返回固定格式的成功响应（code=0）。
func (c *Context) Success(data any) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// Error 返回固定格式的错误响应。
func (c *Context) Error(code int, msg string) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}
