package elasticapm

import (
	"github.com/yyle88/zaplog"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

// CheckApmAgentVersion verifies that APM dependencies use v2
// Returns false if version mismatch is detected
// This check prevents mixing v1 and v2 dependencies which maintain separate tracings
//
// Important context from v1.15.0 to v2.0.0 upgrade:
// - v1: go.elastic.co/apm v1.15.0
// - v2: go.elastic.co/apm/v2 v2.0.0
// The v1 and v2 packages maintain independent defaultTracing instances
//
// CheckApmAgentVersion 验证 APM 依赖是否使用 v2
// 如果检测到版本不匹配则返回 false
// 此检查防止混用 v1 和 v2 依赖，它们维护独立的追踪
//
// 从 v1.15.0 升级到 v2.0.0 的重要背景：
// - v1: go.elastic.co/apm v1.15.0
// - v2: go.elastic.co/apm/v2 v2.0.0
// v1 和 v2 包维护独立的 defaultTracing 实例
func CheckApmAgentVersion(agentVersion string) bool {
	if version := apm.AgentVersion; agentVersion != version {
		zaplog.LOGGER.LOG.Warn("check apm agent versions not match", zap.String("arg_version", agentVersion), zap.String("pkg_version", version))
		return false
	}
	return true
}

// GetApmAgentVersion returns the current APM agent version string
// Version matches the project's go.elastic.co/apm/v2 package version
//
// Examples:
// - If project uses go.elastic.co/apm/v2 v2.6.3, returns "2.6.3"
// - If project uses go.elastic.co/apm/v2 v2.7.0, returns "2.7.0"
// - If project uses go.elastic.co/apm v1.15.0, this returns v2 version (not 1.15.0)
//
// Recommendation: New projects should use go.elastic.co/apm/v2
//
// GetApmAgentVersion 返回当前 APM agent 版本字符串
// 版本与项目的 go.elastic.co/apm/v2 包版本匹配
//
// 示例：
// - 如果项目使用 go.elastic.co/apm/v2 v2.6.3，返回 "2.6.3"
// - 如果项目使用 go.elastic.co/apm/v2 v2.7.0，返回 "2.7.0"
// - 如果项目使用 go.elastic.co/apm v1.15.0，这里返回 v2 版本（不是 1.15.0）
//
// 建议：新项目应使用 go.elastic.co/apm/v2
func GetApmAgentVersion() string {
	return apm.AgentVersion
}
