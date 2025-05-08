package elasticapm

import (
	"github.com/yyle88/zaplog"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

/*
CheckApmAgentVersion 检查是否都是用的 "go.elastic.co/apm/v2" 这个依赖
有些老项目用的还是 "go.elastic.co/apm" 而不是 "go.elastic.co/apm/v2"，这个项目明确只支持新版的 "v2" 依赖
这是我们由 1.15.0 升级为 2.0.0 时发现的问题

	go.elastic.co/apm v1.15.0
	go.elastic.co/apm/module/apmgin v1.15.0
	go.elastic.co/apm/module/apmgopg v1.15.0
	go.elastic.co/apm/module/apmgoredisv8 v1.15.0
	go.elastic.co/apm/module/apmhttp v1.15.0
	go.elastic.co/apm/v2 v2.0.0

这里 v1 和 v2 使用的 defaultTracer 是完全不同的两个变量
*/
func CheckApmAgentVersion(agentVersion string) bool {
	if version := apm.AgentVersion; agentVersion != version {
		zaplog.LOGGER.LOG.Warn("check apm agent versions not match", zap.String("arg_version", agentVersion), zap.String("pkg_version", version))
		return false
	}
	return true
}

// GetApmAgentVersion 获得版本号，只要都是用的 "go.elastic.co/apm/v2" 这个依赖，他就会和主项目的依赖保持相同
// 假如主项目是 go.elastic.co/apm/v2 v2.6.3 就会显示 2.6.3
// 假如主项目是 go.elastic.co/apm/v2 v2.7.0 就会显示 2.7.0
// 但是假如主项目依赖的是 go.elastic.co/apm v1.15.0 这里就不会是 1.15.0 因为这里是 v2 的
// 因此建议新项目都使用 go.elastic.co/apm/v2 这个新版的
func GetApmAgentVersion() string {
	return apm.AgentVersion
}
