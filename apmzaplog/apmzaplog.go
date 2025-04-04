package apmzaplog

import "go.uber.org/zap"

type Log struct {
	sugaredLogger *zap.SugaredLogger
}

func NewLog(sugaredLogger *zap.SugaredLogger) *Log {
	return &Log{sugaredLogger: sugaredLogger}
}

func (o *Log) Debugf(format string, args ...interface{}) {
	o.sugaredLogger.Debugf(format, args...)
}

func (o *Log) Errorf(format string, args ...interface{}) {
	o.sugaredLogger.Errorf(format, args...)
}

func (o *Log) Warningf(format string, args ...interface{}) {
	o.sugaredLogger.Warnf(format, args...)
}
