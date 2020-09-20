package mylog

import (
	"fmt"
	"log"
)

const (
	MYLOG_DEBUG_LEVEL = iota
	MYLOG_ERROR_LEVEL
)

type MyLog struct {
	level  int
	logger *log.Logger
}

func (m MyLog) LogDebug(format string, outpara ...interface{}) {
	if m.level <= MYLOG_DEBUG_LEVEL && m.logger != nil {
		str := fmt.Sprintf("[DEBUG] "+format, outpara...)
		m.logger.Printf(str)
	}
}
func (m MyLog) LogError(format string, outpara ...interface{}) {
	if m.level <= MYLOG_DEBUG_LEVEL && m.logger != nil {
		str := fmt.Sprintf("[ERROR] "+format, outpara...)
		m.logger.Printf(str)
	}
}

var DefaultLogger MyLog

func newLogFile(fileName string, maxAge int) *log.Logger {
	if fileName == "" {
		return nil
	}
	return log.New(newlogger(fileName, maxAge), "", log.LstdFlags)
}

func Init(logFile, logLevel string, maxAge int) error {
	switch logLevel {
	case "debug":
		{
			DefaultLogger.level = MYLOG_DEBUG_LEVEL
		}
	default:
		DefaultLogger.level = MYLOG_ERROR_LEVEL
	}
	DefaultLogger.logger = newLogFile(logFile, maxAge)
	return nil
}

func LogDebug(format string, param ...interface{}) {
	DefaultLogger.LogDebug(format, param...)
}
func LogError(format string, param ...interface{}) {
	DefaultLogger.LogError(format, param...)
}
