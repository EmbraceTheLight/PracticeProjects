package common

type Errno struct {
	Code   int
	Msg    string
	ErrMsg string
}

func (err Errno) Error() string {
	return err.Msg
}

func (err Errno) WithMsg(msg string) Errno {
	err.Msg = err.Msg + "," + msg
	return err
}

func (err Errno) WithError(rawError error) Errno {
	var msg string
	if rawError != nil {
		msg = rawError.Error()
	}
	err.ErrMsg = err.Msg + "," + msg
	return err
}

var (
	Ok            = Errno{Code: 200, Msg: "OK"}
	ServeErr      = Errno{Code: 500, Msg: "Internal Server Error"}
	ParamErr      = Errno{Code: 400, Msg: "Bad Request"}
	AuthErr       = Errno{Code: 401, Msg: "Unauthorized"}
	PermissionErr = Errno{Code: 403, Msg: "Forbidden"}

	DatabaseErr = Errno{Code: 10000, Msg: "Database Error"}
	RedisErr    = Errno{Code: 10001, Msg: "Redis Error"}

	UserNotFoundErr   = Errno{Code: 11001, Msg: "User Not Found"}
	InvalidCaptchaErr = Errno{Code: 11002, Msg: "滑块校验失败, 请重试"}
)
