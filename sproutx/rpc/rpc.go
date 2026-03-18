package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/ink-yht-code/sprout/sproutx/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Config 是 gRPC Server 配置。
type Config struct {
	Enabled bool
	Addr    string
}

// Server 是 gRPC 服务封装。
type Server struct {
	Server *grpc.Server
	addr   string
}

// NewServer 创建 gRPC Server。
func NewServer(cfg Config) *Server {
	if !cfg.Enabled {
		return nil
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)

	return &Server{Server: s, addr: cfg.Addr}
}

// Run 启动 gRPC 服务（阻塞）。
func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	log.Info("gRPC server starting", zap.String("addr", s.addr))
	return s.Server.Serve(lis)
}

// Shutdown 优雅停止 gRPC 服务。
func (s *Server) Shutdown(ctx context.Context) error {
	_ = ctx
	s.Server.GracefulStop()
	return nil
}

func unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	requestID := getRequestID(ctx)
	if requestID != "" {
		ctx = context.WithValue(ctx, "request_id", requestID)
	}

	log.Ctx(ctx).Info("gRPC request",
		zap.String("method", info.FullMethod),
		zap.String("peer", getPeerAddr(ctx)),
	)

	resp, err := handler(ctx, req)
	if err != nil {
		log.Ctx(ctx).Error("gRPC error",
			zap.String("method", info.FullMethod),
			zap.Error(err),
		)
	}

	return resp, err
}

func streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()

	requestID := getRequestID(ctx)
	if requestID != "" {
		ctx = context.WithValue(ctx, "request_id", requestID)
	}

	log.Ctx(ctx).Info("gRPC stream",
		zap.String("method", info.FullMethod),
		zap.String("peer", getPeerAddr(ctx)),
	)

	return handler(srv, ss)
}

func getRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	ids := md.Get("x-request-id")
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}

func getPeerAddr(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	return p.Addr.String()
}

// MapError 将业务错误映射为 gRPC status error。
//
// 若 err 实现 BizCode/BizMsg，则根据 BizCode 的后四位转换为对应 gRPC code。
func MapError(err error) error {
	if err == nil {
		return nil
	}

	type bizError interface {
		BizCode() int
		BizMsg() string
	}

	if biz, ok := err.(bizError); ok {
		code := mapBizCodeToGrpcCode(biz.BizCode())
		return status.Error(code, biz.BizMsg())
	}

	return status.Error(codes.Internal, err.Error())
}

// ErrorDetail 是附加的错误信息结构，用于传递业务码、业务消息与请求 ID。
type ErrorDetail struct {
	BizCode   int32  `json:"biz_code"`
	BizMsg    string `json:"biz_msg"`
	RequestId string `json:"request_id"`
}

func mapBizCodeToGrpcCode(bizCode int) codes.Code {
	suffix := bizCode % 10000
	switch suffix {
	case 1:
		return codes.InvalidArgument
	case 2:
		return codes.Unauthenticated
	case 3:
		return codes.PermissionDenied
	case 4:
		return codes.NotFound
	case 5:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}

type Client struct {
	conn *grpc.ClientConn
}

// NewClient 创建 gRPC Client。
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// Conn 返回底层 grpc.ClientConn。
func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

// Close 关闭连接。
func (c *Client) Close() error {
	return c.conn.Close()
}

// MarshalErrorDetail 将 ErrorDetail 序列化为 JSON 字符串。
func MarshalErrorDetail(detail *ErrorDetail) string {
	data, _ := json.Marshal(detail)
	return string(data)
}

// UnmarshalErrorDetail 将 JSON 字符串反序列化为 ErrorDetail。
func UnmarshalErrorDetail(data string) (*ErrorDetail, error) {
	var detail ErrorDetail
	if err := json.Unmarshal([]byte(data), &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}
