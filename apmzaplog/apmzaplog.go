package apmzaplog

import "go.uber.org/zap"

type Sug struct {
	sug *zap.SugaredLogger
}

func NewSug(sug *zap.SugaredLogger) *Sug {
	return &Sug{sug: sug}
}

func (o *Sug) Debugf(format string, args ...interface{}) {
	o.sug.Debugf(format, args...)
}

func (o *Sug) Errorf(format string, args ...interface{}) {
	o.sug.Errorf(format, args...)
}

func (o *Sug) Warningf(format string, args ...interface{}) {
	o.sug.Warnf(format, args...)
}
