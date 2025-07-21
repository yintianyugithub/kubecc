package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"greet/core/nacos"
)

type Config struct {
	zrpc.RpcServerConf
	ClusterName string `json:"-"`
	Secret      string `json:",optional"`
	Host        string `json:",optional"`
	Port        int64  `json:",optional"`
}

// GetCnf 从Nacos获取配置
func (c *Config) GetCnf(nc *nacos.Config) string {
	ct, err := nc.Inc.GetConfig(vo.ConfigParam{
		DataId: nc.AddRpcName,
		Group:  nc.AddRpcName,
	})

	if err != nil {
		panic(err)
	}

	if ct == "" {
		panic(fmt.Errorf("获取配置失败, DataId: %s, Group: %s", nc.AddRpcName, nc.AddRpcName))
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
		return
	}

	// 取消服务注册
	//_, err = nc.IClient.DeregisterInstance(vo.DeregisterInstanceParam{
	//	Ip:          c.Host,
	//	Port:        8080,
	//	ServiceName: "adm-rpc",
	//	GroupName:   "adm-rpc",
	//})

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

}
