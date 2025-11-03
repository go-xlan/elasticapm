// Package apmzaplog provides zap logging integration with Elastic APM
// Implements apm.Logging interface to route APM logs through zap
// Uses frame skipping to ensure correct source location in logs
//
// apmzaplog 提供 zap 日志与 Elastic APM 的集成
// 实现 apm.Logging 接口将 APM 日志路由到 zap
// 使用帧跳过确保日志中的源位置正确
package apmzaplog

import (
	"github.com/yyle88/zaplog"
)

// Log implements apm.Logging interface using zap logging
// Routes APM logs to the application's zap logging
// Uses Skip(1) to report correct source information
//
// Log 使用 zap 日志实现 apm.Logging 接口
// 将 APM 日志路由到应用的 zap 日志
// 使用 Skip(1) 报告正确的源信息
type Log struct{}

// NewLog creates a new zap-based APM logging
// Returns instance that can be set as APM's default logging
//
// NewLog 创建新的基于 zap 的 APM 日志
// 返回可设置成 APM 默认日志的实例
func NewLog() *Log {
	return &Log{}
}

// Debugf logs a debug message from APM using zap
// Skips one frame to show APM's source location
//
// Debugf 使用 zap 记录来自 APM 的调试消息
// 跳过一帧以显示 APM 的源位置
func (o *Log) Debugf(format string, args ...interface{}) {
	zaplog.ZAPS.Skip(1).SUG.Debugf(format, args...)
}

// Errorf logs an error message from APM using zap
// Skips one frame to show APM's source location
//
// Errorf 使用 zap 记录来自 APM 的错误消息
// 跳过一帧以显示 APM 的源位置
func (o *Log) Errorf(format string, args ...interface{}) {
	zaplog.ZAPS.Skip(1).SUG.Errorf(format, args...)
}

// Warningf logs a warning message from APM using zap
// Skips one frame to show APM's source location
//
// Warningf 使用 zap 记录来自 APM 的警告消息
// 跳过一帧以显示 APM 的源位置
func (o *Log) Warningf(format string, args ...interface{}) {
	zaplog.ZAPS.Skip(1).SUG.Warnf(format, args...)
}
