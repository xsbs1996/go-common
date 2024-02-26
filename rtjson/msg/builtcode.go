package msg

const (
	Ok                     = 0  // 成功
	HttpErr                = 1  // http错误
	UnknownErr             = 2  // 未知错误
	UnsetErr               = 3  // 未设置错误
	CodeRangeErr           = 4  // 错误码区间错误
	UnknownBusinessTypeErr = 5  // 未知的服务类型
	CodeDuplicationErr     = 6  // 错误码重复
	RequestParameterErr    = 7  // 请求参数错误
	ServerInternalError    = 8  // 服务器内部错误
	FrequentOperations     = 9  // 操作频繁
	UserNotLoggedIn        = 10 // 用户未登录
	MemberCodeErr          = 11 // 验证码错误
	PhoneNumberErr         = 12 // 手机号错误
	EmailErr               = 13 // 邮箱错误
	InsufficientBalanceErr = 14 // 余额不足
	ForbiddenLogin         = 15 // 禁止登录
	ForbiddenGame          = 16 // 禁止游戏
	ForbiddenTrade         = 17 // 禁止交易

)

var BuiltMessage = map[int]string{
	Ok:                     "ok",
	HttpErr:                "http error",
	UnknownErr:             "unknown error",
	UnsetErr:               "unset error message",
	CodeRangeErr:           "error code range is incorrect",
	UnknownBusinessTypeErr: "unknown service type error",
	CodeDuplicationErr:     "error code duplication",
	RequestParameterErr:    "request parameter error",
	ServerInternalError:    "the network is busy. please try again later",
	FrequentOperations:     "frequent operation",
	UserNotLoggedIn:        "not logged in",
	MemberCodeErr:          "verification code error",
	PhoneNumberErr:         "phone number error",
	EmailErr:               "email error",
	InsufficientBalanceErr: "insufficient balance",
	ForbiddenLogin:         "Account Disable",
	ForbiddenGame:          "forbidden game",
	ForbiddenTrade:         "forbidden trade",
}
