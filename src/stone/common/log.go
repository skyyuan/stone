package common

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

var (
	// Logger debug logger
	Logger echo.Logger
)

// LoggerInit init globel logger
func LoggerInit(e *echo.Echo, debug bool) {
	Logger = e.Logger
	if debug {
		Logger.SetLevel(log.DEBUG)
	} else {
		Logger.SetLevel(log.INFO)
	}
}
