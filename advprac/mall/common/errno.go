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

func (err Errno) IsOk() bool {
	return err.Code == 200
}
func (err Errno) NotOk() bool {
	return err.IsOk() == false
}

var (
	Ok            = Errno{Code: 200, Msg: "OK"}
	ServeErr      = Errno{Code: 500, Msg: "Internal Server Error"}
	ParamErr      = Errno{Code: 400, Msg: "Bad Request"}
	AuthErr       = Errno{Code: 401, Msg: "Unauthorized"}
	PermissionErr = Errno{Code: 403, Msg: "Forbidden"}

	DatabaseErr = Errno{Code: 10000, Msg: "Database Error"}
	RedisErr    = Errno{Code: 10001, Msg: "Redis Error"}

	UserNotFoundErr    = Errno{Code: 11001, Msg: "用户不存在"}
	InvalidCaptchaErr  = Errno{Code: 11002, Msg: "滑块校验失败, 请重试"}
	InvalidPasswordErr = Errno{Code: 11003, Msg: "用户不存在或密码错误"}
	PasswordErrLimit   = Errno{Code: 11004, Msg: "用户名或密码错误次数超过限定次数, 请十分钟后再试"}
	JwtCreateTokenErr  = Errno{Code: 11005, Msg: "登录失败"}
	UserNeedLogin      = Errno{Code: 11006, Msg: "登录过期, 请重新登录"}
)
