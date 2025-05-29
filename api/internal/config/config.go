package config

import (
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"greet/core/nacos"
)

func (c *Config) HotLoadCnf(nc *nacos.Config) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("Recovered from hot load config panic: %v", r)
		}
	}()
	err := nc.Inc.ListenConfig(vo.ConfigParam{
		DataId: nc.Name,
		Group:  nc.Name,
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

	Secret string `json:",optional"`
}
