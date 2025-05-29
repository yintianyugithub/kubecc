package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/logx"
)

type PointParam struct {
	NamespaceId string
	ServerName  string
	GroupName   string
	DataId      string
	Inc         config_client.IConfigClient

	c *constant.ClientConfig
	s []constant.ServerConfig
}

// Register 服务注册、发现
func (p *PointParam) Register() {
	p.c = &constant.ClientConfig{
		// 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		NamespaceId:         p.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	p.s = []constant.ServerConfig{
		{
			IpAddr:      "127.0.0.1",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	// 创建服务发现客户端
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  p.c,
			ServerConfigs: p.s,
		},
	)

	if err != nil {
		panic(err)
	}

	// 服务发现-注册
	_, err = namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: p.ServerName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: fmt.Sprintf("cluster-%s", p.GroupName), // 默认值DEFAULT
		GroupName:   p.GroupName,                            // 默认值DEFAULT_GROUP
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

	p.Inc = p.newConfig()
}

// GetCnf 从Nacos获取配置
func (p *PointParam) GetCnf() string {
	ct, err := p.Inc.GetConfig(vo.ConfigParam{
		DataId: p.DataId,
		Group:  p.GroupName,
	})

	if err != nil {
		panic(err)
	}

	if ct == "" {
		panic(fmt.Errorf("获取配置失败, DataId: %s, Group: %s", p.DataId, p.GroupName))
	}

	return ct
}

// HotLoadCnf todo 热加载配置 data已更新，但服务配置未更新

// newConfig 初始化nacos配置
func (p *PointParam) newConfig() config_client.IConfigClient {
	// 创建动态配置客户端
	cnfClt, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  p.c,
			ServerConfigs: p.s,
		},
	)

	if err != nil {
		panic(err)
	}

	return cnfClt
}
