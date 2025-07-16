package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"greet/api/internal/config"
	"greet/api/internal/middleware"
)

type ServiceContext struct {
	Config           *config.Config
	SvcGvcMiddleware rest.Middleware
	JwtMiddleware    rest.Middleware
}

func NewServiceContext(c *config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		SvcGvcMiddleware: middleware.NewSvcGvcMiddleware().Handle,
		JwtMiddleware:    middleware.NewJwtMiddleware().Handle,
	}
}
