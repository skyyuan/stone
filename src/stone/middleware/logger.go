package middleware

import (
	"math"
	"time"

	"stone/common"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// Logger is the logrus logger handler
func Logger(log *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			path := req.RequestURI
			if err = next(c); err != nil {
				c.Error(err)
			}

			stop := time.Since(start)
			latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000.0))
			statusCode := res.Status
			clientIP := c.RealIP()
			// clientUserAgent := req.UserAgent()
			// referer := req.Referer()
			// hostname, err := os.Hostname()
			// if err != nil {
			// 	hostname = "unknow"
			// }
			dataLength := res.Size
			if dataLength < 0 {
				dataLength = 0
			}

			headers := req.Header

			requestID := headers.Get(echo.HeaderXRequestID)

			errcode, errmsg, _ := parseCtxError(c, err)

			body := []interface{}{statusCode, latency, clientIP, req.Method, path, dataLength, errcode, errmsg}

			entry := logrus.NewEntry(log).WithFields(
				logrus.Fields{
					"requestID": requestID,
					"body":      body,
				})

			entry.Info("")
			return
		}
	}
}

func parseCtxError(c echo.Context, err error) (errcode string, errmsg string, panicStack string) {
	errcode = "0"

	// 优先处理panic错误
	if panicErr := c.Get(common.ErrorCtx); panicErr != nil {
		v := panicErr.(*common.BizError)
		errcode = v.Code
		errmsg = v.Msg
		panicStack = v.Stack
	} else if err != nil {
		switch v := err.(type) {
		case *common.BizError:
			errcode = v.Code
			errmsg = v.Msg
			panicStack = v.Stack
		default:
			errmsg = v.Error()
		}
	}

	return
}
