package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"greet/core/nacos"
)

type Config struct {
	rest.RestConf
	ClusterName string `json:"-"`
	Secret      string `json:",optional"`

	AddRpc zrpc.RpcClientConf `json:",optional"`
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

// Register 服务注册
func (c *Config) Register(nc *nacos.Config) {
	c.ClusterName = fmt.Sprintf("cluster-%s", c.Name)

	_, err := nc.IClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          c.Host,
		Port:        uint64(c.Port),
		ServiceName: c.Name,
		ClusterName: c.ClusterName,
		GroupName:   c.Name,
	})
	if err != nil {
		logx.Error(err)
	}

	c.hotLoadCnf(nc)
}

func (c *Config) hotLoadCnf(nc *nacos.Config) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("Recovered from hot load config panic: %v", r)
		}
	}()

	// 监听配置变化
	if err := nc.Inc.ListenConfig(vo.ConfigParam{
		DataId: c.Name,
		Group:  c.Name,
		OnChange: func(namespace, group, dataId, data string) {
			logx.Infof("api config changed, namespace:%s group:%s dataId:%s data:%s", namespace, group, dataId, data)
			if err := conf.LoadFromYamlBytes([]byte(data), c); err != nil {
				panic(err)
			}
		},
	}); err != nil {
		logx.Error(err)
	}

	// 监听服务变化
	if err := nc.IClient.Subscribe(&vo.SubscribeParam{
		ServiceName: c.Name,
		Clusters:    []string{c.ClusterName},
		GroupName:   c.Name,
		SubscribeCallback: func(services []model.Instance, err error) {
			logx.Infof("api service changed：%s", util.ToJsonString(services))
		},
	}); err != nil {
		logx.Error(err)
	}

}

// GetService 获取服务信息
func (c *Config) GetService(nc *nacos.Config) (model.Service, error) {
	svcInfo, err := nc.IClient.GetService(vo.GetServiceParam{
		Clusters:    []string{c.ClusterName},
		ServiceName: c.Name,
		GroupName:   c.Name,
	})
	if err != nil {
		logx.Error(err)
		return model.Service{}, err
	}

	return svcInfo, nil
}

// SelectOneHealthyInstance 获取一个健康实例(负载均衡策略：加权随机轮询)
func (c *Config) SelectOneHealthyInstance(nc *nacos.Config) (*model.Instance, error) {
	healthIns, err := nc.IClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		Clusters:    []string{c.ClusterName},
		ServiceName: c.Name,
		GroupName:   c.Name,
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	return healthIns, nil
}
