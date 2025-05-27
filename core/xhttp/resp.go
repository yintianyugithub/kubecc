package xhttp

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/net/context"
	"greet/core/xerr"
	"net/http"
	"reflect"
	"strings"
)

type Body struct {
	RequestId string       `json:"request_id"` // 请求ID
	Code      xerr.ErrCode `json:"code"`
	Msg       string       `json:"msg"`
	Data      interface{}  `json:"data"`
}

// ResponseParamsInvalid 通用参数异常响应
func ResponseParamsInvalid(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	path := "none"
	if r.URL != nil {
		path = r.URL.Path
	}

	//todo FIXME:埋点
	_ = path
	//stats.ApiRespCodeCounter.Inc(path, strconv.FormatInt(int64(xerr.ErrCodeParamsInvalid), 10))

	httpx.OkJsonCtx(ctx, w, Body{
		RequestId: trace.TraceIDFromContext(ctx),
		Code:      xerr.ErrCodeParamsInvalid,
		Msg:       fmt.Sprintf("参数错误: %s", strings.Replace(err.Error(), "\"", "'", -1)),
	})
}

// Response 包装 code msg 统一返回结构体
func Response(ctx context.Context, w http.ResponseWriter, r *http.Request, resp interface{}, err error) {
	code := xerr.ErrCodeNone
	msg := "ok"

	if err != nil {
		// 兜底错误
		code = xerr.ErrCodeUnknown
		msg = err.Error()

		// 业务错误
		var e *xerr.AppErr
		if v := errors.As(err, &e); v {
			code = e.Code()
		}
	}

	// todo: FIXME:埋点

	httpx.OkJsonCtx(ctx, w, &Body{
		RequestId: trace.TraceIDFromContext(ctx),
		Code:      code,
		Msg:       msg,
		Data:      normalizeNil(resp),
	})
}

func normalizeNil(i interface{}) interface{} {
	if i == nil {
		return nil
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Func:
		if v.IsNil() {
			return nil
		}
	default:
		return i
	}
	return i
}
