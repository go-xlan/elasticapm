package elasticapm

import (
	"github.com/go-xlan/elasticapm/apmzaplog"
	"github.com/yyle88/erero"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

// Initialize 初始化全局的APM，使用的是 "go.elastic.co/apm/v2" 而非旧版的 ""go.elastic.co/apm" 依赖
func Initialize(cfg *Config) error {
	if err := InitializeWithOptions(cfg, NewEnvOption()); err != nil {
		return erero.Wro(err)
	}
	zaplog.LOG.Debug("Initialize apm success")
	//设置日志
	apm.DefaultTracer().SetLogger(apmzaplog.NewLog())
	return nil
}

func InitializeWithOptions(cfg *Config, evo *EnvOption, setEnvs ...func()) error {
	zaplog.SUG.Info("Initialize apm cfg=" + neatjsons.S(cfg))
	zaplog.SUG.Info("Initialize apm evo=" + neatjsons.S(evo))

	if cfg.ServerUrl == "" && len(cfg.ServerUrls) == 0 {
		return erero.New("APM server URL is missing")
	}

	cfg.SetEnv(evo)

	for _, setFunc := range setEnvs {
		setFunc() //设置环境变量
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
