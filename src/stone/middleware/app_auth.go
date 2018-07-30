package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo"
	"stone/common"
	"stone/common/auth"
)

// AppAuth 鉴权检查
func AppAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			header := c.Request().Header
			appkey := header.Get("X-App-Key")
			if appkey == "" {
				return common.BizError9001
			}
			nonce := header.Get("X-Nonce")
			if nonce == "" {
				return common.BizError9001
			}
			timestamp := header.Get("X-Timestamp")
			if timestamp == "" {
				return common.BizError9001
			}
			checksum := header.Get("X-App-Checksum")
			if checksum == "" {
				return common.BizError9001
			}
			unixts := time.Now().Unix()
			ntimestamp, _ := strconv.Atoi(timestamp)
			if unixts-int64(ntimestamp) > 300.0 {
				return common.BizError9000
			}

			check, err := auth.AppVerify(appkey, nonce, timestamp, checksum)
			if err != nil || !check {
				return common.BizError9001
			}
			// auth ok, call next handle
			return next(c)
		}
	}
}
