package main

import (
	"flag"
	"fmt"

	{{.importPackages}}
)

var (
    configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")
    nc = nacos.PointParam{
        NamespaceId: "",
        ServerName:  "",
        GroupName:   "",
        DataId:      "",
    }
)

func main() {
	flag.Parse()

	c := &config.Config{}
	conf.MustLoad(*configFile, c)

    nc.Register()
	if err := conf.LoadFromYamlBytes([]byte(nc.GetCnf()), c); err != nil {
		panic(err)
		return
	}
	go c.HotLoadCnf(nc.DataId, nc.GroupName)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
