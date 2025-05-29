package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strconv"
	"strings"
)

type Config struct {
	Name string `json:",optional"`
	Host string `json:",optional"`
	Port int    `json:",optional"`
	Mode string `json:",optional"`

	Nacos   Nacos                       `json:",optional"`
	Inc     config_client.IConfigClient `json:",optional"`
	C       *constant.ClientConfig      `json:",optional"`
	S       []constant.ServerConfig     `json:",optional"`
	IClient naming_client.INamingClient `json:",optional"`
}

type Nacos struct {
	Hosts  []string `json:",hosts"`
	Key    string   `json:",key"`
	Scheme string   `json:",scheme"`
}

func (c *Config) Init() {
	port, err := strconv.Atoi(strings.Split(c.Nacos.Hosts[0], ":")[1])
	if err != nil {
		panic(err)
	}
	ip := strings.Split(c.Nacos.Hosts[0], ":")[0]

	c.C = &constant.ClientConfig{
		// 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		NamespaceId:         c.Mode,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	c.S = []constant.ServerConfig{
		{
			IpAddr:      ip,
			ContextPath: "/nacos",
			Port:        uint64(port),
			Scheme:      c.Nacos.Scheme,
		},
	}

	// 创建服务发现客户端
	c.IClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  c.C,
			ServerConfigs: c.S,
		},
	)

	// nacos自身-注册
	_, err = c.IClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: "nacos",
		ClusterName: "cluster-nacos",
		GroupName:   "nacos",
	})
	if err != nil {
		panic(err)
	}

	// 服务注册
	_, err = c.IClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          c.Host,
		Port:        uint64(c.Port),
		ServiceName: c.Name,
		ClusterName: "cluster-" + c.Name,
		GroupName:   c.Name,
	})
	if err != nil {
		panic(err)
	}

	c.Inc = c.newConfig()

	go c.Sub()
}

// GetCnf 从Nacos获取配置
func (c *Config) GetCnf() string {
	ct, err := c.Inc.GetConfig(vo.ConfigParam{
		DataId: c.Name,
		Group:  c.Name,
	})

	if err != nil {
		panic(err)
	}

	if ct == "" {
		panic(fmt.Errorf("获取配置失败, DataId: %s, Group: %s", c.Name, c.Name))
	}

	return ct
}

// newConfig 初始化nacos配置
func (c *Config) newConfig() config_client.IConfigClient {
	// 创建动态配置客户端
	cnfClt, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  c.C,
			ServerConfigs: c.S,
		},
	)

	if err != nil {
		panic(err)
	}

	return cnfClt
}

func (c *Config) Sub() {
	// todo fix-me // 服务发现-注销、获取服务信息、获取所有实例列表、获取实例列表、获取一个健康实例(加权轮询)
	// todo 取消服务监听、获取服务名列表
	// 监听服务变化
	if err := c.IClient.Subscribe(&vo.SubscribeParam{
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
	svcInfo, err := c.IClient.GetService(vo.GetServiceParam{
		Clusters:    []string{fmt.Sprintf("cluster-%s", c.Name)},
		ServiceName: c.Name,
		GroupName:   c.Name,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(util.ToJsonString(svcInfo))
}
