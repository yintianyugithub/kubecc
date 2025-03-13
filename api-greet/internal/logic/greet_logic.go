package logic

import (
	"context"
	"greet/service/add/pb/add"

	"greet/api-greet/internal/svc"
	"greet/api-greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GreetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGreetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GreetLogic {
	return &GreetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GreetLogic) Greet(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.RpcAdd.Add(l.ctx, &add.AddReq{})
	return
}
