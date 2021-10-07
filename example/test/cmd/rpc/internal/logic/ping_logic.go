package logic

import (
	"context"
	"time"

	"git.zc0901.com/go/god/example/test/cmd/rpc/internal/svc"
	"git.zc0901.com/go/god/example/test/cmd/rpc/test"

	"git.zc0901.com/go/god/lib/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(req *test.PingReq) (*test.PingReply, error) {
	time.Sleep(5 * time.Second)
	return &test.PingReply{
		Pong: req.Name,
	}, nil
}
