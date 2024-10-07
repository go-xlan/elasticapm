package elasticapm

import (
	"github.com/yyle88/zaplog"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

/*
CheckApmAgentVersion 检查你的apm包的版本和我的是否相同
在这里，同一个项目，肯定是相等的
但是假如保持这里的代码长期保持不变，当将来再写新项目时，这时假如第三方库由 apm "2.0.0" 升级为 "2.0.1" 就会导致问题
毕竟调用 utils 的 init 函数赋值给 defaultTracer 的是 2.0.0 包里面的全局变量，而你在新项目里使用的是 2.0.1 包里的全局变量
因此我们建议你在引用这个 utils 包的项目里，在项目开头都写这样的断言
或者比较完以后不panic而是发出告警、日志等信息
因为两个包使用的版本不相同时，监控肯定是无效的
这是我们由 1.15.0 升级为 2.0.0 时发现的问题，相对来说还是挺坑的

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

// GetApmAgentVersion 获得版本号
// 假如你用的我的包，我的包 go.mod 里面引用的 apm 是 v2.0.0 的，这里就会返回 v2.0.0 版本号字符串
// 假如你项目里直接用到 apm 包，则你需要检查你的包和我的包版本是否相同
// 假如不同，逻辑不通
func GetApmAgentVersion() string {
	return apm.AgentVersion
}
