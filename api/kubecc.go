package main

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"greet/api/internal/config"
	"greet/api/internal/handler"
	"greet/api/internal/svc"
	"greet/core/nacos"
)

func main() {

	c := &config.Config{}
	nc := nacos.Init()
	if err := conf.LoadFromYamlBytes([]byte(c.GetCnf(nc)), c); err != nil {
		panic(err)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	go c.Register(nc)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}
