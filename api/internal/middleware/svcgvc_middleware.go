package middleware

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type SvcGvcMiddleware struct {
}

func NewSvcGvcMiddleware() *SvcGvcMiddleware {
	return &SvcGvcMiddleware{}
}

func (m *SvcGvcMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logx.Infof("Request received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		next(w, r)
	}
}
