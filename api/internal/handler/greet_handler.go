package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"greet/api/internal/logic"
	"greet/api/internal/svc"
	"greet/api/internal/types"
	"greet/core/xhttp"
)

func GreetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.ResponseParamsInvalid(r.Context(), w, r, err)
			return
		}

		l := logic.NewGreetLogic(r.Context(), svcCtx)
		resp, err := l.Greet(&req)
		xhttp.Response(r.Context(), w, r, resp, err)
	}
}
