package logic

import (
	"context"

	"git.zc0901.com/go/god/example/graceful/dns/api/internal/svc"
	"git.zc0901.com/go/god/example/graceful/dns/api/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type GracefulLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGracefulLogic(ctx context.Context, svcCtx *svc.ServiceContext) GracefulLogic {
	return GracefulLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GracefulLogic) Graceful() (*types.Response, error) {
	// todo: add your logic here and delete this line

	return &types.Response{}, nil
}
