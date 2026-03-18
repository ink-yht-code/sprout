package rpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestNewServer_Disabled(t *testing.T) {
	cfg := Config{Enabled: false, Addr: ":9090"}
	s := NewServer(cfg)
	if s != nil {
		t.Fatalf("expected nil server when disabled")
	}
}

func TestNewServer_Enabled(t *testing.T) {
	cfg := Config{Enabled: true, Addr: ":0"}
	s := NewServer(cfg)
	if s == nil {
		t.Fatalf("expected non-nil server when enabled")
	}
	if s.Server == nil {
		t.Fatalf("expected grpc server to be initialized")
	}
}

func TestServer_Shutdown(t *testing.T) {
	cfg := Config{Enabled: true, Addr: ":0"}
	s := NewServer(cfg)
	if s == nil {
		t.Fatalf("expected non-nil server")
	}

	ctx := context.Background()
	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}

func TestInterceptorsExist(t *testing.T) {
	ctx := context.Background()
	_, _ = unaryInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, func(ctx context.Context, req any) (any, error) {
		return nil, nil
	})

	_ = streamInterceptor(nil, &dummyStream{ctx: ctx}, &grpc.StreamServerInfo{FullMethod: "/test"}, func(srv any, stream grpc.ServerStream) error {
		return nil
	})
}

type dummyStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (d *dummyStream) Context() context.Context {
	return d.ctx
}
