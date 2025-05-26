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

var (
	nc = nacos.PointParam{
		NamespaceId: "b7ed918d-562a-4d15-85d7-d46dd7852262",
		ServerName:  "api",
		GroupName:   "api",
		DataId:      "1",
	}
	c config.Config
)

func main() {
	nc.Register()

	go nc.HotLoadCnf()

	if err := conf.LoadFromYamlBytes([]byte(nc.GetCnf()), &c); err != nil {
		panic(err)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, svc.NewServiceContext(c))

	fmt.Printf("🚀 Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}
