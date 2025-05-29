package main

import (
	"flag"
	"fmt"
	"greet/core/nacos"

	"greet/api/internal/config"
	"greet/api/internal/handler"
	"greet/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	{{.importPackages}}
)

var (
    // todo FIXME: 修改名称（configFile、nc）
    configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")
    nc = &nacos.PointParam{
		ServerName:  "{{.serviceName}}",
		GroupName:   "{{.serviceName}}",
		DataId:      "{{.serviceName}}",
    }
)

func main() {
	flag.Parse()

	c := &config.Config{}
	conf.MustLoad(*configFile, c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	c.Register(nc)
	go c.HotLoadCnf(nc)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
