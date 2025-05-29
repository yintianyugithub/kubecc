package logic

import (
	"context"
	"greet/core/xerr"

	"greet/api/internal/svc"
	"greet/api/internal/types"

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

	err = xerr.Error(xerr.ErrCodeParamsInvalid, "参数错误")
	resp = &types.Response{
		Message:  l.svcCtx.Config.Secret,
		Message1: "1",
	}
	return
}
