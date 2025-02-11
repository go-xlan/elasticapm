package elasticapm

import (
	"fmt"
	"os"
	"strings"

	"github.com/yyle88/must"
)

type Config struct {
	Environment    string
	ServerUrls     []string
	ServerUrl      string
	ApiKey         string `json:"-"`
	SecretToken    string `json:"-"`
	ServiceName    string
	ServiceVersion string
	NodeName       string
	ServerCertPath string
	SkipShortSpans bool
}

func (cfg *Config) SetEnv(evo *EnvOption) {
	evo.SetEnv("ELASTIC_APM_ACTIVE", fmt.Sprint(true))
	evo.SetEnv("ELASTIC_APM_ENVIRONMENT", cfg.Environment)

	evo.SetEnv("ELASTIC_APM_SERVER_URLS", strings.Join(cfg.ServerUrls, ","))
	evo.SetEnv("ELASTIC_APM_SERVER_URL", cfg.ServerUrl)

	evo.SetEnv("ELASTIC_APM_API_KEY", cfg.ApiKey)
	evo.SetEnv("ELASTIC_APM_SECRET_TOKEN", cfg.SecretToken)
	if cfg.NodeName != "" {
		evo.SetEnv("ELASTIC_APM_SERVICE_NODE_NAME", cfg.NodeName)
	}
	if cfg.ServerCertPath != "" {
		evo.SetEnv("ELASTIC_APM_SERVER_CERT", cfg.ServerCertPath)
	}
	if !cfg.SkipShortSpans { //假如不忽略，就配置0ms（但是没有用，有些组件会把小于1ms的截掉精度，再当成不大于0ms的，因此postgres 的 gorm里面小于1ms的还是会被忽略掉的）
		evo.SetEnv("ELASTIC_APM_SPAN_STACK_TRACE_MIN_DURATION", "0ms")
	}
}

type EnvOption struct {
	Override bool //假如设为false就会使用系统环境变量的值
}

func NewEnvOption() *EnvOption {
	return &EnvOption{
		Override: true,
	}
}

func (op *EnvOption) SetEnv(name, value string) {
	if op.Override {
		must.Done(os.Setenv(name, value))
		return
	}

	if os.Getenv(name) == "" {
		must.Done(os.Setenv(name, value))
		return
	}
}
