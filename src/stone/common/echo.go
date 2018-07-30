package common

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// EchoHTTPErrorHandler is a HTTP error handler. It sends a JSON response
// with status code.
func EchoHTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			code = http.StatusOK
			msg  interface{}
			rmsg string
		)

		errcode := ErrorCode9999
		if he, ok := err.(*echo.HTTPError); ok {
			msg = he.Message
		} else if be, ok := err.(*BizError); ok {
			errcode = be.Code
			msg = be.TMsg(c)
		} else if e.Debug {
			msg = err.Error()
		} else {
			msg = http.StatusText(http.StatusInternalServerError)
		}
		// 处理错误信息
		if v, ok := msg.(string); ok {
			rmsg = v
		} else {
			rmsg = fmt.Sprintf("%s", msg)
		}

		if !c.Response().Committed {
			if c.Request().Method == echo.HEAD { // Issue #608
				if err := c.NoContent(code); err != nil {
					goto ERROR
				}
			} else {
				// 统一封装返回值
				if err := c.JSON(code, ErrorReturns(errcode, rmsg)); err != nil {
					goto ERROR
				}
			}
		}
	ERROR:
		e.Logger.Error(err)
	}
}

// GetAcceptLanguage Get Accept-Language from request header
func GetAcceptLanguage(c echo.Context) string {
	return c.Request().Header.Get("Accept-Language")
}

// EchoInit echo init
func EchoInit(e *echo.Echo, conf Config) {
	e.Debug = conf.Debug()
	LoggerInit(e, conf.Debug())
	e.HTTPErrorHandler = EchoHTTPErrorHandler(e)
}
