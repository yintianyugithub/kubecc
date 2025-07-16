package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"greet/core/nacos"
)

type Config struct {
	rest.RestConf

	Secret string `json:",optional"`
}

// Register 服务注册
func (c *Config) Register(nc *nacos.Config) {
	_, err := nc.IClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          c.Host,
		Port:        uint64(c.Port),
		ServiceName: c.Name,
		ClusterName: "cluster-" + c.Name,
		GroupName:   c.Name,
	})
	if err != nil {
		panic(err)
	}

	c.hotLoadCnf(nc)
}

func (c *Config) hotLoadCnf(nc *nacos.Config) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("Recovered from hot load config panic: %v", r)
		}
	}()
	err := nc.Inc.ListenConfig(vo.ConfigParam{
		DataId: c.Name,
		Group:  c.Name,
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

// 监听服务变化
func (c *Config) sub(nc *nacos.Config) {
	if err := nc.IClient.Subscribe(&vo.SubscribeParam{
		ServiceName: c.Name,
		Clusters:    []string{fmt.Sprintf("cluster-%s", c.Name)},
		GroupName:   c.Name,
		SubscribeCallback: func(services []model.Instance, err error) {
			fmt.Println("服务变化回调：", util.ToJsonString(services))
		},
	}); err != nil {
		panic(err)
	}

	// 获取服务信息
	svcInfo, err := nc.IClient.GetService(vo.GetServiceParam{
		Clusters:    []string{fmt.Sprintf("cluster-%s", c.Name)},
		ServiceName: c.Name,
		GroupName:   c.Name,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(util.ToJsonString(svcInfo))
}

// GetCnf 从Nacos获取配置
func (c *Config) GetCnf(nc *nacos.Config) string {
	ct, err := nc.Inc.GetConfig(vo.ConfigParam{
		DataId: nc.ApiName,
		Group:  nc.ApiName,
	})

	if err != nil {
		panic(err)
	}

	if ct == "" {
		panic(fmt.Errorf("获取配置失败, DataId: %s, Group: %s", c.Name, c.Name))
	}

	return ct
}
