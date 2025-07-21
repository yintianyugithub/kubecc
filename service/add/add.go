package main

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"greet/core/nacos"
	"greet/service/add/internal/config"
	"greet/service/add/internal/server"
	"greet/service/add/internal/svc"
	"greet/service/add/pb/add"
)

func main() {
	c := &config.Config{}
	nc := nacos.Init()

	if err := conf.LoadFromYamlBytes([]byte(c.GetCnf(nc)), &c); err != nil {
		panic(err)
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		add.RegisterAdderServer(grpcServer, server.NewAdderServer(svc.NewServiceContext(c)))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	threading.GoSafe(func() {
		c.Register(nc)
	})

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)

	s.Start()
}
