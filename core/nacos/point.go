package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type PointParam struct {
	NamespaceId string
	ServerName  string
	GroupName   string
	DataId      string
	Inc         config_client.IConfigClient

	C *constant.ClientConfig
	S []constant.ServerConfig
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

// NewConfig 初始化nacos配置
func (p *PointParam) NewConfig() config_client.IConfigClient {
	// 创建动态配置客户端
	cnfClt, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  p.C,
			ServerConfigs: p.S,
		},
	)

	if err != nil {
		panic(err)
	}

	return cnfClt
}
