package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
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
)

func main() {
	nc.Register()
	c := &config.Config{}
	hot(c)

	if err := conf.LoadFromYamlBytes([]byte(nc.GetCnf()), c); err != nil {
		panic(err)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, svc.NewServiceContext(c))

	fmt.Printf("🚀 Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}

func hot(c *config.Config) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("Recovered from hot load config panic: %v", r)
		}
	}()
	confi := nacos.NewConfig()
	err := confi.ListenConfig(vo.ConfigParam{
		DataId: nc.DataId,
		Group:  nc.GroupName,
		OnChange: func(namespace, group, dataId, data string) {
			logx.Info(namespace, group, dataId, data)
			err := conf.LoadFromYamlBytes([]byte(data), c)
			if err != nil {
				logx.Errorf("Failed to load config: %v", err)
				return
			}
		},
	})
	if err != nil {
		panic(err)
	}
}
