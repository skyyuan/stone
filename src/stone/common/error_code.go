package common

var (
	// ErrorCtx context key for error
	ErrorCtx = "ERR_CTX"

	// ErrorCode0 正常值
	ErrorCode0 = "0"

	// ErrorCode9999 系统异常
	ErrorCode9999 = "9999"

	// BizError1001 调用参数异常
	BizError1001 = NewBizError("1001")
	// BizError9000 请求鉴权超时
	BizError9000 = NewBizError("9000")
	// BizError9001 请求鉴权非法
	BizError9001 = NewBizError("9001")
)
