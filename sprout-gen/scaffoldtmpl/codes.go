package scaffoldtmpl

// CodesTmpl 错误码模板
var CodesTmpl = `package errs

// ServiceID 服务 ID
// 错误码格式：ServiceID * 10000 + BizCode
// 示例：ServiceID=01，BizCode=1 → 错误码 = 010001
const ServiceID = {{.ServiceID}}

// 业务码定义
const (
	CodeSuccess       = ServiceID*10000 + 0
	CodeInvalidParam  = ServiceID*10000 + 1
	CodeUnauthorized  = ServiceID*10000 + 2
	CodeForbidden     = ServiceID*10000 + 3
	CodeNotFound      = ServiceID*10000 + 4
	CodeConflict      = ServiceID*10000 + 5
	CodeInternalError = ServiceID*10000 + 9999
	
	// TODO: 添加业务特定错误码
	// 示例：
	// Code{{.NameUpper}}NotFound = ServiceID*10000 + 10
	// Code{{.NameUpper}}NameDuplicate = ServiceID*10000 + 11
)

// 预定义错误
var (
	ErrSuccess       = NewBizError(CodeSuccess, "success")
	ErrInvalidParam  = NewBizError(CodeInvalidParam, "invalid parameter")
	ErrUnauthorized  = NewBizError(CodeUnauthorized, "unauthorized")
	ErrForbidden     = NewBizError(CodeForbidden, "forbidden")
	ErrNotFound      = NewBizError(CodeNotFound, "not found")
	ErrConflict      = NewBizError(CodeConflict, "conflict")
	ErrInternalError = NewBizError(CodeInternalError, "internal error")
)
`
