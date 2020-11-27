package logic

import (
	"context"

	"git.zc0901.com/go/god/example/graceful/dns/rpc/graceful"
	"git.zc0901.com/go/god/example/graceful/dns/rpc/internal/svc"

	"git.zc0901.com/go/god/lib/logx"
)

type GraceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGraceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GraceLogic {
	return &GraceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GraceLogic) Grace(in *graceful.Request) (*graceful.Response, error) {
	// todo: add your logic here and delete this line

	return &graceful.Response{}, nil
}
