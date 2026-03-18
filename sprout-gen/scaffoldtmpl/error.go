package scaffoldtmpl

// ErrorTmpl 错误类型模板
var ErrorTmpl = `package errs

import "fmt"

// BizError 业务错误
type BizError struct {
	Code  int
	Msg   string
	Cause error
}

func (e *BizError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Msg, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}

func (e *BizError) Unwrap() error {
	return e.Cause
}

// NewBizError 创建业务错误
func NewBizError(code int, msg string, cause ...error) *BizError {
	var c error
	if len(cause) > 0 {
		c = cause[0]
	}
	return &BizError{Code: code, Msg: msg, Cause: c}
}

// InvalidParam 参数错误
func InvalidParam(msg string, cause ...error) *BizError {
	return NewBizError(CodeInvalidParam, msg, cause...)
}

// Unauthorized 未授权
func Unauthorized(msg string, cause ...error) *BizError {
	return NewBizError(CodeUnauthorized, msg, cause...)
}

// Forbidden 无权限
func Forbidden(msg string, cause ...error) *BizError {
	return NewBizError(CodeForbidden, msg, cause...)
}

// NotFound 未找到
func NotFound(msg string, cause ...error) *BizError {
	return NewBizError(CodeNotFound, msg, cause...)
}

// Conflict 冲突
func Conflict(msg string, cause ...error) *BizError {
	return NewBizError(CodeConflict, msg, cause...)
}

// InternalError 内部错误
func InternalError(msg string, cause ...error) *BizError {
	return NewBizError(CodeInternalError, msg, cause...)
}
`
