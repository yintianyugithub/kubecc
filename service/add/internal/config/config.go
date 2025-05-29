package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	Nacos Nacos  `json:",optional"`
	Mode  string `json:",optional"` // namespace id
}

type Nacos struct {
	Hosts  []string `json:",optional"`
	Scheme string   `json:",optional"`
	Key    string   `json:",optional"`
}
