package xerr

type AppErr struct {
	//typ  ErrType
	code ErrCode
	msg  string
}

func (e *AppErr) Error() string {
	return e.msg
}

func (e *AppErr) Code() ErrCode {
	return e.code
}

func Error(code ErrCode, msg string) *AppErr {
	return &AppErr{
		code: code,
		msg:  msg,
	}
}
