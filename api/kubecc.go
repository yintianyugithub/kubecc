package main

import (
	"flag"
	"fmt"
	"greet/api/internal/config"
	"greet/api/internal/handler"
	"greet/api/internal/svc"
	"greet/core/nacos"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var (
	// todo FIXME: kubecc-api.yaml -> kubecc-nacos.yaml
	configFile = flag.String("f", "etc/kubecc-api.yaml", "the config file")
)

func main() {
	flag.Parse()

	nc := &nacos.Config{}

	conf.MustLoad(*configFile, nc)

	nc.Init()

	c := &config.Config{}
	if err := conf.LoadFromYamlBytes([]byte(nc.GetCnf()), c); err != nil {
		panic(err)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	go c.HotLoadCnf(nc)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
