package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// AccessLog 表示一次请求的访问日志结构。
type AccessLog struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Query    string `json:"query"`
	IP       string `json:"ip"`
	UserID   string `json:"user_id"`
	ReqBody  string `json:"req_body"`
	RespBody string `json:"resp_body"`
	Status   int    `json:"status"`
	Duration int64  `json:"duration"`
	Error    string `json:"error"`
}

// AccessLogFunc 是访问日志回调函数。
type AccessLogFunc func(log *AccessLog)

// AccessLogBuilder 用于构建访问日志中间件。
type AccessLogBuilder struct {
	logFunc       AccessLogFunc
	logReqBody    bool
	logRespBody   bool
	maxBodyLength int
}

// NewAccessLogBuilder 创建一个 AccessLogBuilder。
func NewAccessLogBuilder(logFunc AccessLogFunc) *AccessLogBuilder {
	return &AccessLogBuilder{
		logFunc:       logFunc,
		logReqBody:    false,
		logRespBody:   false,
		maxBodyLength: 1024,
	}
}

// WithReqBody 配置是否记录请求体。
func (b *AccessLogBuilder) WithReqBody(log bool) *AccessLogBuilder {
	b.logReqBody = log
	return b
}

// WithRespBody 配置是否记录响应体。
func (b *AccessLogBuilder) WithRespBody(log bool) *AccessLogBuilder {
	b.logRespBody = log
	return b
}

// WithMaxBodyLength 配置记录 body 的最大长度。
func (b *AccessLogBuilder) WithMaxBodyLength(length int) *AccessLogBuilder {
	b.maxBodyLength = length
	return b
}

// Build 构建访问日志中间件。
func (b *AccessLogBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log := &AccessLog{
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Query:  c.Request.URL.RawQuery,
			IP:     c.ClientIP(),
		}

		if userId, exists := c.Get("user_id"); exists {
			if uid, ok := userId.(string); ok {
				log.UserID = uid
			}
		}

		if b.logReqBody && c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > b.maxBodyLength {
				log.ReqBody = string(bodyBytes[:b.maxBodyLength]) + "...(truncated)"
			} else {
				log.ReqBody = string(bodyBytes)
			}
		}

		if b.logRespBody {
			writer := &accessLogResponseWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
			}
			c.Writer = writer

			c.Next()

			respBody := writer.body.String()
			if len(respBody) > b.maxBodyLength {
				log.RespBody = respBody[:b.maxBodyLength] + "...(truncated)"
			} else {
				log.RespBody = respBody
			}
		} else {
			c.Next()
		}

		log.Status = c.Writer.Status()
		log.Duration = time.Since(start).Milliseconds()

		if len(c.Errors) > 0 {
			log.Error = c.Errors.String()
		}

		if b.logFunc != nil {
			b.logFunc(log)
		}
	}
}

type accessLogResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *accessLogResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}
