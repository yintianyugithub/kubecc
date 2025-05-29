package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"greet/core/nacos"
	"greet/service/add/internal/config"
	"greet/service/add/internal/server"
	"greet/service/add/internal/svc"
	"greet/service/add/pb/add"
)

var (
	configFile = flag.String("f", "etc/add.yaml", "the config file")
	nc         = nacos.PointParam{
		NamespaceId: "b7ed918d-562a-4d15-85d7-d46dd7852262",
		ServerName:  "add",
		GroupName:   "api",
		DataId:      "2",
	}
)

func main() {
	flag.Parse()

	c := &config.Config{}

	conf.MustLoad(*configFile, c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		add.RegisterAdderServer(grpcServer, server.NewAdderServer(svc.NewServiceContext(c)))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
