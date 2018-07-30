package common

import (
	"fmt"

	"github.com/labstack/echo"
)

// BizError 业务错误
type BizError struct {
	Code  string
	Msg   string
	Stack string
}

// NewBizError 生成一个BizError
// code: 错误码
// msg: 可变参数，分别为：错误信息、panic堆栈（不必须）
func NewBizError(code string, msg ...string) *BizError {
	var emsg string
	if len(msg) > 0 {
		emsg = msg[0]
	}
	var stack string
	if len(msg) > 1 {
		stack = msg[1]
	}
	return &BizError{code, emsg, stack}
}

// 实现error接口
func (bize BizError) Error() string {
	return fmt.Sprintf("BizError:%s/%s", bize.Code, bize.Msg)
}

// TMsg return translated message
func (bize *BizError) TMsg(c echo.Context) string {
	if bize.Msg != "" {
		return bize.Msg
	}
	return TLang(c)("errcode_" + bize.Code)
}
