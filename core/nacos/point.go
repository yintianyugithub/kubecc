package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Hosts   []string                    `json:",hosts"`
	Key     string                      `json:",key"`
	Scheme  string                      `json:",scheme"`
	Mode    string                      `json:",optional"`
	Inc     config_client.IConfigClient `json:",optional"`
	C       *constant.ClientConfig      `json:",optional"`
	S       []constant.ServerConfig     `json:",optional"`
	IClient naming_client.INamingClient `json:",optional"`

	// 服务配置文件名
	ApiName string `json:",optional"`
}

var (
	cnf  *Config
	once sync.Once
)

func Init() *Config {

	once.Do(func() {

		conf.MustLoad("./core/nacos/config.yaml", &cnf)

		port, err := strconv.Atoi(strings.Split(cnf.Hosts[0], ":")[1])
		if err != nil {
			panic(err)
		}
		ip := strings.Split(cnf.Hosts[0], ":")[0]

		cnf.C = &constant.ClientConfig{
			// 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
			NamespaceId:         cnf.Mode,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			LogLevel:            "debug",
		}

		cnf.S = []constant.ServerConfig{
			{
				IpAddr:      ip,
				ContextPath: "/nacos",
				Port:        uint64(port),
				Scheme:      cnf.Scheme,
			},
		}

		// 创建服务发现客户端
		cnf.IClient, err = clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  cnf.C,
				ServerConfigs: cnf.S,
			},
		)

		// 注册nacos
		if _, err = cnf.IClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          ip,
			Port:        uint64(port),
			ServiceName: "nacos",
			ClusterName: "cluster-nacos",
			GroupName:   "nacos",
		}); err != nil {
			panic(err)
		}

		// 初始化nacos配置
		if cnf.Inc, err = clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  cnf.C,
				ServerConfigs: cnf.S,
			},
		); err != nil {
			panic(err)
		}
	})

	return cnf
}
