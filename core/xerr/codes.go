package xerr

//go:generate stringer -type=ErrCode -linecomment -output strings.go
type ErrCode int64

// 错误码枚举，请勿随意修改错误码枚举值
const (
	ErrCodeNone          ErrCode = 0     // OK
	ErrCodeServerError   ErrCode = 500   // 服务异常
	ErrCodeUnknown       ErrCode = 10000 // 未知错误
	ErrCodeParamsInvalid ErrCode = 10001 // 参数错误
	ErrCodeTokenInvalid  ErrCode = 10002 // Token 无效
	ErrCodeAuthFailed    ErrCode = 10003 // 认证失败
)
