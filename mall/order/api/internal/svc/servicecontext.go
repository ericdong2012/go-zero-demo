package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-demo/mall/order/api/internal/config"
	"go-zero-demo/mall/user/rpc/userclient"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
