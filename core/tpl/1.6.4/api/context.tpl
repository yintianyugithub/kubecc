package svc

import (
	{{.configImport}}
)

type ServiceContext struct {
	Config *{{.config}}
	{{.middleware}}
}

func NewServiceContext(c {{.config}}) *ServiceContext {
	return &ServiceContext{
		Config: c,
		{{.middlewareAssignment}}
	}
}
