package server_interceptors

import (
	"context"
	"fmt"
	"git.zc0901.com/go/god/lib/logx"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

func init() {
	logx.Disable()
}

func TestUnaryCrashInterceptor(t *testing.T) {
	interceptor := UnaryCrashInterceptor()
	_, err := interceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("mock panic")
	})
	assert.NotNil(t, err)
	fmt.Print(err)
}

func TestStreamCrashInterceptor(t *testing.T) {
	err := StreamCrashInterceptor(nil, nil, nil, func(srv interface{}, stream grpc.ServerStream) error {
		panic("mock panic")
	})
	assert.NotNil(t, err)
	fmt.Print(err)
}
