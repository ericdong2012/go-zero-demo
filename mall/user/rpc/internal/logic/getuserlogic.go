package logic

import (
	"context"
	"go-zero-demo/mall/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-demo/mall/user/rpc/internal/svc"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *userclient.IdRequest) (*userclient.UserResponse, error) {
	return &userclient.UserResponse{
		Id:   in.Id,
		Name: "test1",
	}, nil
}
