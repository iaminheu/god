package logic

import (
	"context"

	"git.zc0901.com/go/god/example/test/cmd/rpc/test"

	"git.zc0901.com/go/god/example/test/cmd/api/internal/svc"
	"git.zc0901.com/go/god/example/test/cmd/api/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) PingLogic {
	return PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingLogic) Ping(req types.PingReq) (*types.PingReply, error) {
	pingReply, err := l.svcCtx.TestRPC.Ping(l.ctx, &test.PingReq{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &types.PingReply{
		Pong: pingReply.Pong,
	}, nil
}
