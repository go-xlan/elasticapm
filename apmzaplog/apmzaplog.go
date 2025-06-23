package apmzaplog

import (
	"github.com/yyle88/zaplog"
)

type Log struct{}

func NewLog() *Log {
	return &Log{}
}

func (o *Log) Debugf(format string, args ...interface{}) {
	zaplog.ZAPS.Skip(1).SUG.Debugf(format, args...)
}

func (o *Log) Errorf(format string, args ...interface{}) {
	zaplog.ZAPS.Skip(1).SUG.Errorf(format, args...)
}

func (o *Log) Warningf(format string, args ...interface{}) {
	zaplog.ZAPS.Skip(1).SUG.Warnf(format, args...)
}
