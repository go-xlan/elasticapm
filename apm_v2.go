package elasticapm

import (
	"github.com/go-xlan/elasticapm/apmzaplog"
	"github.com/yyle88/erero"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

// INIT 初始化全局的APM，需要校验 agentVersion 是因为有可能系统里是 2.0.0 而新项目用的是 2.0.1 的，或者总之就是出现版本不匹配的情况
func INIT(cfg *Config) error {
	var evo = NewEnvOption()

	zaplog.SUG.Info("init apm cfg=" + neatjsons.S(cfg))
	zaplog.SUG.Info("init apm evo=" + neatjsons.S(evo))

	if err := INIT2(cfg, evo); err != nil {
		return erero.Wro(err)
	}

	zaplog.LOG.Debug("init apm success")
	//设置日志
	apm.DefaultTracer().SetLogger(apmzaplog.NewSug(zaplog.ZAPS.P1.SUG))
	return nil
}

func INIT2(cfg *Config, evo *EnvOption, setEnvs ...func()) error {
	if cfg.ServerUrx == "" && len(cfg.ServerUrls) == 0 {
		zaplog.LOG.Debug("no apm urx. not init apm")
		return erero.New("no apm urx")
	}

	cfg.SetEnv(evo)

	for _, set := range setEnvs {
		set() //设置环境变量
	}

	zaplog.LOG.Debug("init apm", zap.String("service_name", cfg.ServiceName), zap.String("agent_version", apm.AgentVersion))

	atc, err := apm.NewTracer(cfg.ServiceName, cfg.ServiceVersion)
	if err != nil {
		return erero.Wro(err)
	}
	apm.SetDefaultTracer(atc)

	zaplog.LOG.Debug("init apm success")
	return nil
}
