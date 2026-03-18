package core

import "errors"

// ErrNoResponse 表示当前 Handler 不需要写回响应（用于 SSE/流式输出等场景）。
// ErrUnauthorized 表示未授权（会被包装器转换为 401）。
// ErrSessionNotFound 表示会话不存在。
// ErrSessionExpired 表示会话已过期。
// ErrInvalidToken 表示 JWT 无效。
var (
	ErrNoResponse = errors.New("不需要返回响应")

	ErrUnauthorized = errors.New("未授权")

	ErrSessionNotFound = errors.New("会话不存在")

	ErrSessionExpired = errors.New("会话已过期")

	ErrInvalidToken = errors.New("无效的令牌")
)
