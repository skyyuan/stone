package common

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// RequestLog 返回一个logrus.Logger对象
func RequestLog(reqlog string) *logrus.Logger {

	reqlogger := logrus.New()
	if reqlog != "" {
		lfs := lfshook.NewHook(lfshook.PathMap{
			logrus.InfoLevel: reqlog,
		})
		lfs.SetFormatter(&RequestFormatter{})
		reqlogger.Hooks.Add(lfs)
		// disable standard output
		reqlogger.Out = ioutil.Discard
	} else {
		reqlogger.Formatter = &RequestFormatter{}
	}

	return reqlogger
}

// RequestFormatter requestFormatter
type RequestFormatter struct {
}

var logTemplate = []string{"timestamp", "requestID", "module", "funcName", "level", "body"}

const (
	logTimeFormat = "2006-01-02 15:04:05.000"
	logSpliter    = "\t\u0001"
)

// Format renders a single log entry
func (f *RequestFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var res []string

	for _, item := range logTemplate {
		switch item {
		case "timestamp":
			res = append(res, entry.Time.Format(logTimeFormat))
		case "level":
			res = append(res, "->")
			res = append(res, strings.ToUpper(entry.Level.String()))
		default:
			res = append(res, plainText(entry.Data[item]))
		}
	}
	res = append(res, entry.Message, "\n")
	return []byte(strings.Join(res, logSpliter)), nil
}

func plainText(item interface{}) string {
	switch v := item.(type) {
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(item.(int64), 10)
	case string:
		return v
	case []interface{}:
		var res []string
		for _, vi := range v {
			res = append(res, plainText(vi))
		}
		return strings.Join(res, logSpliter)
	default:
		return ""
	}
}
