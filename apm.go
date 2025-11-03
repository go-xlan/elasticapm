// Package elasticapm provides a simple and elegant wrapping of Elastic APM v2
// Offers simple initialization, configuration management, and integration with zap logging
// Supports distributed tracing across HTTP and gRPC services
//
// elasticapm 提供简洁优雅的 Elastic APM v2 封装
// 提供便捷的初始化、配置管理和 zap 日志集成
// 支持 HTTP 和 gRPC 服务的分布式追踪
package elasticapm

import (
	"github.com/go-xlan/elasticapm/apmzaplog"
	"github.com/yyle88/erero"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

// Initialize sets up the APM tracing with given config and default options
// Uses "go.elastic.co/apm/v2" instead of the old "go.elastic.co/apm"
// Sets up zap logging integration with APM logs
//
// Initialize 使用给定配置和默认选项初始化 APM 追踪
// 使用 "go.elastic.co/apm/v2" 而非旧版 "go.elastic.co/apm"
// 设置 zap 日志与 APM 日志的集成
func Initialize(cfg *Config) error {
	if err := InitializeWithOptions(cfg, NewEnvOption()); err != nil {
		return erero.Wro(err)
	}
	zaplog.LOG.Debug("Initialize apm success")
	// Set zap logging integration with APM tracing
	// 设置 zap 日志与 APM 追踪的集成
	apm.DefaultTracer().SetLogger(apmzaplog.NewLog())
	return nil
}

// InitializeWithOptions sets up APM with custom environment options
// Allows fine-grained controls on environment variable handling
// Extra setup functions can be passed when custom configuration is needed
//
// InitializeWithOptions 使用自定义环境选项初始化 APM
// 允许对环境变量处理进行精细控制
// 可传递额外的设置函数进行自定义配置
func InitializeWithOptions(cfg *Config, evo *EnvOption, setEnvs ...func()) error {
	zaplog.SUG.Info("Initialize apm cfg=" + neatjsons.S(cfg))
	zaplog.SUG.Info("Initialize apm evo=" + neatjsons.S(evo))

	if cfg.ServerUrl == "" && len(cfg.ServerUrls) == 0 {
		return erero.New("APM server URL is missing")
	}

	cfg.SetEnv(evo)

	// Apply extra environment setup functions
	// 应用额外的环境设置函数
	for _, setFunc := range setEnvs {
		setFunc()
	}

	zaplog.LOG.Debug("Initialize apm", zap.String("service_name", cfg.ServiceName), zap.String("agent_version", apm.AgentVersion))

	tracer, err := apm.NewTracer(cfg.ServiceName, cfg.ServiceVersion)
	if err != nil {
		return erero.Wro(err)
	}
	apm.SetDefaultTracer(tracer)

	zaplog.LOG.Debug("Initialize apm success")
	return nil
}
