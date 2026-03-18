package scaffoldtmpl

// ResultTmpl result 模板（响应结果封装）
var ResultTmpl = `package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 统一响应结果
type Result struct {
	Code int         ` + "`" + `json:"code"` + "`" + `
	Msg  string      ` + "`" + `json:"msg"` + "`" + `
	Data interface{} ` + "`" + `json:"data,omitempty"` + "`" + `
}

// Success 成功响应
func Success(data interface{}) Result {
	return Result{Code: 0, Msg: "success", Data: data}
}

// Error 错误响应
func Error(code int, msg string) Result {
	return Result{Code: code, Msg: msg}
}

// JSON 输出 JSON
func JSON(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Result{Code: code, Msg: "", Data: data})
}
`
