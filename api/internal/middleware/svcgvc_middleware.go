package middleware

import (
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/zeromicro/go-zero/core/logx"
	"greet/api/internal/config"
	"greet/core/nacos"
	"net/http"
)

type SvcGvcMiddleware struct {
	c *config.Config
}

func NewSvcGvcMiddleware(c *config.Config) *SvcGvcMiddleware {
	return &SvcGvcMiddleware{
		c: c,
	}
}

func (m *SvcGvcMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logx.Infof("Request received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		ins, err := m.c.SelectOneHealthyInstance(nacos.Init())
		if err != nil {
			logx.Errorf("Failed to select a healthy instance: %v", err)
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		// todo 咋转发请求到选中的健康实例
		logx.Infof("Forwarding request to instance: %s", util.ToJsonString(ins))

		next(w, r)
	}
}
