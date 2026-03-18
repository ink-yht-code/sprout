package core

// CodeSuccess 等业务码常量用于统一返回结果的 Code 字段。
const (
	// CodeSuccess 表示成功。
	CodeSuccess = 0

	// CodeWarning 表示警告（业务可继续）。
	CodeWarning = 1

	// CodeError 表示通用错误。
	CodeError = 2

	// CodeInvalidParam 表示参数错误。
	CodeInvalidParam = 10000

	// CodeInternalError 表示系统内部错误。
	CodeInternalError = 20000

	// CodeUnauthorized 表示未授权。
	CodeUnauthorized = 20001

	// CodeForbidden 表示无权限。
	CodeForbidden = 20003

	// CodeNotFound 表示资源不存在。
	CodeNotFound = 20004

	// CodeConflict 表示资源冲突。
	CodeConflict = 20005

	// CodeTooManyRequests 表示请求过于频繁。
	CodeTooManyRequests = 20006

	// CodeServiceUnavailable 表示服务不可用。
	CodeServiceUnavailable = 20007
)

// CodeMessage 是默认业务码到中文消息的映射表。
var CodeMessage = map[int]string{
	CodeSuccess:            "成功",
	CodeWarning:            "警告",
	CodeError:              "错误",
	CodeInvalidParam:       "参数错误",
	CodeInternalError:      "系统繁忙",
	CodeUnauthorized:       "未授权",
	CodeForbidden:          "没有权限",
	CodeNotFound:           "资源不存在",
	CodeConflict:           "资源冲突",
	CodeTooManyRequests:    "请求过于频繁",
	CodeServiceUnavailable: "服务不可用",
}

// GetCodeMessage 根据业务码返回默认的中文提示。
func GetCodeMessage(code int) string {
	if msg, ok := CodeMessage[code]; ok {
		return msg
	}

	switch {
	case code == 0:
		return "成功"
	case code == 1:
		return "警告"
	case code == 2:
		return "错误"
	default:
		return "未知错误"
	}
}

// Success 构造成功响应。
func Success(msg string, data any) Result {
	if msg == "" {
		msg = "成功"
	}
	return Result{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	}
}

// Warning 构造警告响应。
func Warning(msg string, data any) Result {
	if msg == "" {
		msg = "警告"
	}
	return Result{
		Code: CodeWarning,
		Msg:  msg,
		Data: data,
	}
}

// Error 构造通用错误响应。
func Error(msg string) Result {
	if msg == "" {
		msg = "错误"
	}
	return Result{
		Code: CodeError,
		Msg:  msg,
		Data: nil,
	}
}

// ErrorWithCode 构造指定业务码的错误响应。
func ErrorWithCode(code int, msg string) Result {
	if msg == "" {
		msg = GetCodeMessage(code)
	}
	return Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

// InvalidParam 构造参数错误响应。
func InvalidParam(msg string) Result {
	if msg == "" {
		msg = GetCodeMessage(CodeInvalidParam)
	}
	return Result{Code: CodeInvalidParam, Msg: msg, Data: nil}
}

// InternalError 构造系统内部错误响应。
func InternalError() Result {
	return Result{Code: CodeInternalError, Msg: GetCodeMessage(CodeInternalError), Data: nil}
}

// Unauthorized 构造未授权响应。
func Unauthorized() Result {
	return Result{Code: CodeUnauthorized, Msg: GetCodeMessage(CodeUnauthorized), Data: nil}
}

// Forbidden 构造无权限响应。
func Forbidden() Result {
	return Result{Code: CodeForbidden, Msg: GetCodeMessage(CodeForbidden), Data: nil}
}

// NotFound 构造资源不存在响应。
func NotFound(msg string) Result {
	if msg == "" {
		msg = GetCodeMessage(CodeNotFound)
	}
	return Result{Code: CodeNotFound, Msg: msg, Data: nil}
}

// Conflict 构造资源冲突响应。
func Conflict(msg string) Result {
	if msg == "" {
		msg = GetCodeMessage(CodeConflict)
	}
	return Result{Code: CodeConflict, Msg: msg, Data: nil}
}

// TooManyRequests 构造请求过于频繁响应。
func TooManyRequests() Result {
	return Result{Code: CodeTooManyRequests, Msg: GetCodeMessage(CodeTooManyRequests), Data: nil}
}

// ServiceUnavailable 构造服务不可用响应。
func ServiceUnavailable() Result {
	return Result{Code: CodeServiceUnavailable, Msg: GetCodeMessage(CodeServiceUnavailable), Data: nil}
}
