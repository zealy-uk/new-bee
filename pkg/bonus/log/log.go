package log

import (
	"github.com/newswarm-lab/new-bee/pkg/logging"
)

var logger logging.Logger

//Init ...
func Init(initlogger logging.Logger) {
	logger = initlogger
}

//Trace ...
func Trace(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

//Debug ...
func Debug(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

//Info ...
func Info(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

//Warn ...
func Warn(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

//Error ...
func Error(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

//Panic ...
func Panic(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
