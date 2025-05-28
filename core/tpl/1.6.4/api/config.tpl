package config

import (
    "github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
    {{.authImport}}
)

type Config struct {
	rest.RestConf
	{{.auth}}
	{{.jwtTrans}}

	Nc         config_client.IConfigClient `conf:"-"`
	Nacos      Nacos `conf:",optional"`
}

type Nacos struct {
    Hosts []string `conf:"hosts"`
    Key   string `conf:"key"`
    Scheme string `conf:"scheme"`
}

func (c *Config) HotLoadCnf(dataId, groupId string) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("Recovered from hot load config panic: %v", r)
		}
	}()
	err := c.Nc.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  groupId,
		OnChange: func(namespace, group, dataId, data string) {
			logx.Info(namespace, group, dataId, data)
			if err := conf.LoadFromYamlBytes([]byte(data), c); err != nil {
				panic(err)
			}
		},
	})
	if err != nil {
		panic(err)
	}
}