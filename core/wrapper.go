package core

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ink-yht-code/sprout/context"
	"github.com/ink-yht-code/sprout/session"
)

func handleResponse(c *gin.Context, res Result, err error, userId string) {
	if errors.Is(err, ErrNoResponse) {
		slog.Debug("不需要响应", slog.Any("err", err))
		return
	}

	if errors.Is(err, ErrUnauthorized) {
		slog.Debug("未授权", slog.Any("err", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err != nil {
		logFields := []any{
			slog.String("path", c.Request.URL.Path),
			slog.Any("err", err),
		}
		if userId != "" {
			logFields = append(logFields, slog.String("user_id", userId))
		}
		slog.Error("执行业务逻辑失败", logFields...)

		if res.Code == 0 {
			res = InternalError()
		} else {
			res.Data = nil
			if res.Msg == "" {
				res.Msg = GetCodeMessage(res.Code)
			}
		}
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// W 将业务函数包装为 gin.HandlerFunc。
//
// 适用于无需请求体绑定的场景。业务函数返回 (Result, error)：
// - error != nil 时会记录日志并返回统一错误响应
// - 返回 ErrNoResponse 时不会写回响应
// - 返回 ErrUnauthorized 时会直接返回 401
func W(fn func(ctx *context.Context) (Result, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &context.Context{Context: c}
		res, err := fn(ctx)
		handleResponse(c, res, err, "")
	}
}

// B 将业务函数包装为 gin.HandlerFunc，并自动进行请求参数绑定。
//
// Req 为请求结构体类型（建议带 validate tag）。绑定失败会直接返回 400。
func B[Req any](fn func(ctx *context.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &context.Context{Context: c}

		var req Req
		if err := c.ShouldBind(&req); err != nil {
			slog.Debug("绑定参数失败",
				slog.String("path", c.Request.URL.Path),
				slog.Any("err", err))
			c.JSON(http.StatusBadRequest, InvalidParam("参数错误"))
			return
		}

		res, err := fn(ctx, req)
		handleResponse(c, res, err, "")
	}
}

// S 将业务函数包装为 gin.HandlerFunc，并自动获取 Session。
//
// 获取 Session 失败会直接返回 401。
func S(fn func(ctx *context.Context, sess session.Session) (Result, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &context.Context{Context: c}

		sess, err := session.Get(ctx)
		if err != nil {
			slog.Debug("获取 Session 失败",
				slog.String("path", c.Request.URL.Path),
				slog.Any("err", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		res, err := fn(ctx, sess)
		handleResponse(c, res, err, sess.Claims().UserId)
	}
}

// BS 将业务函数包装为 gin.HandlerFunc，同时支持参数绑定和 Session 获取。
//
// 获取 Session 失败返回 401；绑定失败返回 400。
func BS[Req any](fn func(ctx *context.Context, req Req, sess session.Session) (Result, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &context.Context{Context: c}

		sess, err := session.Get(ctx)
		if err != nil {
			slog.Debug("获取 Session 失败",
				slog.String("path", c.Request.URL.Path),
				slog.Any("err", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var req Req
		if err := c.ShouldBind(&req); err != nil {
			slog.Debug("绑定参数失败",
				slog.String("path", c.Request.URL.Path),
				slog.String("user_id", sess.Claims().UserId),
				slog.Any("err", err))
			c.JSON(http.StatusBadRequest, InvalidParam("参数错误"))
			return
		}

		res, err := fn(ctx, req, sess)
		handleResponse(c, res, err, sess.Claims().UserId)
	}
}
