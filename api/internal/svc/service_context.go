package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"greet/api/internal/config"
	"greet/api/internal/middleware"
	"greet/service/add/adder"
	"time"
)

type ServiceContext struct {
	Config           *config.Config
	SvcGvcMiddleware rest.Middleware
	JwtMiddleware    rest.Middleware

	AddRpc adder.Adder
}

func NewServiceContext(c *config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		SvcGvcMiddleware: middleware.NewSvcGvcMiddleware(c).Handle,
		JwtMiddleware:    middleware.NewJwtMiddleware(c).Handle,

		AddRpc: adder.NewAdder(zrpc.MustNewClient(c.AddRpc, zrpc.WithTimeout(time.Duration(c.AddRpc.Timeout)*time.Millisecond))),
	}
}
