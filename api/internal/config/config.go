package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"greet/core/nacos"
	"strconv"
	"strings"
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

// Register 服务注册、发现
func (c *Config) Register(p *nacos.PointParam) {
	p.NamespaceId = c.Mode

	port, err := strconv.Atoi(strings.Split(c.Nacos.Hosts[0], ":")[1])
	if err != nil {
		panic(err)
	}
	ip := strings.Split(c.Nacos.Hosts[0], ":")[0]

	p.C = &constant.ClientConfig{
		// 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		NamespaceId:         p.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	p.S = []constant.ServerConfig{
		{
			IpAddr:      ip,
			ContextPath: "/nacos",
			Port:        uint64(port),
			Scheme:      c.Nacos.Scheme,
		},
	}

	// 创建服务发现客户端
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  p.C,
			ServerConfigs: p.S,
		},
	)

	if err != nil {
		panic(err)
	}

	// 服务发现-注册
	_, err = namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: p.ServerName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: fmt.Sprintf("cluster-%s", p.GroupName),
		GroupName:   p.GroupName,
	})

	if err != nil {
		panic(err)
	}

	// todo fix-me // 服务发现-注销、获取服务信息、获取所有实例列表、获取实例列表、获取一个健康实例(加权轮询)
	// todo 取消服务监听、获取服务名列表
	// 服务注册-监听服务变化
	if err := namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: p.ServerName,
		Clusters:    []string{fmt.Sprintf("cluster-%s", p.GroupName)},
		GroupName:   p.GroupName,
		SubscribeCallback: func(services []model.Instance, err error) {
			logx.Infof("服务变化回调: %s, 错误: %+v", util.ToJsonString(services), err)
		},
	}); err != nil {
		panic(err)
	}

	p.Inc = p.NewConfig()

	if err := conf.LoadFromYamlBytes([]byte(p.GetCnf()), c); err != nil {
		panic(err)
	}
}

type Config struct {
	rest.RestConf

	Nacos  Nacos  `json:",optional"`
	Secret string `json:",optional"`
}

type Nacos struct {
	Hosts  []string `json:"hosts"`
	Key    string   `json:"key"`
	Scheme string   `json:"scheme"`
}
