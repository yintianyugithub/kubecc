package {{.PkgName}}

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"greet/core/xhttp"
	{{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
        if err := httpx.Parse(r, &req); err != nil {
			xhttp.ResponseParamsInvalid(r.Context(), w, r, err)
			return
		}

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		{{if .HasResp}}xhttp.Response(r.Context(), w, r, resp, err){{else}}xhttp.Response(r.Context(), w, r, nil, err){{end}}
	}
}
