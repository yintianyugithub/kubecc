package config

import (
    "github.com/nacos-group/nacos-sdk-go/v2/vo"
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/rest"
    "greet/core/nacos"
    {{.authImport}}
)

func (c *Config) HotLoadCnf(nc *nacos.PointParam) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("Recovered from hot load config panic: %v", r)
		}
	}()
	err := nc.Inc.ListenConfig(vo.ConfigParam{
		DataId: nc.DataId,
		Group:  nc.GroupName,
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

type Config struct {
	rest.RestConf
	{{.auth}}
	{{.jwtTrans}}

	Nacos  Nacos  `json:",optional"`
	Secret string `json:",optional"`
}

type Nacos struct {
	Hosts  []string `json:"hosts"`
	Key    string   `json:"key"`
	Scheme string   `json:"scheme"`
}
