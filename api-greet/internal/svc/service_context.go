package svc

import (
	"greet/api-greet/internal/config"
	"greet/service/add/adder"
)

type ServiceContext struct {
	Config config.Config
	RpcAdd adder.Adder
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
