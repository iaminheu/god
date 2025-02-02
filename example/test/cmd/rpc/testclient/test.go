// Code generated by god. DO NOT EDIT!
// Source: test.proto

//go:generate mockgen -destination ./test_mock.go -package testclient -source $GOFILE

package testclient

import (
	"context"

	"git.zc0901.com/go/god/example/test/cmd/rpc/test"

	"git.zc0901.com/go/god/rpc"
)

type (
	PingReply = test.PingReply
	PingReq   = test.PingReq

	Test interface {
		Ping(ctx context.Context, req *PingReq) (*PingReply, error)
	}

	defaultTest struct {
		cli rpc.Client
	}
)

func NewTest(cli rpc.Client) Test {
	return &defaultTest{
		cli: cli,
	}
}

func (m *defaultTest) Ping(ctx context.Context, req *PingReq) (*PingReply, error) {
	client := test.NewTestClient(m.cli.Conn())
	return client.Ping(ctx, req)
}
