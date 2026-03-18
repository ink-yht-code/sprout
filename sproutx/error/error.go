package error

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BizError 表示业务错误接口。
//
// BizCode 为业务码（通常为 ServiceID*10000 + suffix）；BizMsg 为面向用户的提示信息。
type BizError interface {
	BizCode() int
	BizMsg() string
	Error() string
}

// MapToHTTP 将错误映射为 HTTP 响应。
//
// 若 err 为 BizError，则根据 BizCode 后四位映射到对应 HTTP status。
func MapToHTTP(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var biz BizError
	if errors.As(err, &biz) {
		status, resp := mapBizError(biz)
		c.JSON(status, resp)
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    0,
		"message": "internal error",
	})
}

func mapBizError(biz BizError) (int, gin.H) {
	bizCode := biz.BizCode()
	suffix := bizCode % 10000

	var status int
	switch suffix {
	case 1:
		status = http.StatusBadRequest
	case 2:
		status = http.StatusUnauthorized
	case 3:
		status = http.StatusForbidden
	case 4:
		status = http.StatusNotFound
	case 5:
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
	}

	return status, gin.H{
		"code":    bizCode,
		"message": biz.BizMsg(),
	}
}

// Handler 返回一个错误处理中间件。
//
// 当请求链路中出现 gin.Errors 时，将第一个错误映射为 HTTP 响应。
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			MapToHTTP(c, err)
		}
	}
}
