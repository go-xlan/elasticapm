package elasticapm

import (
	"fmt"
	"os"
	"strings"

	"github.com/yyle88/must"
)

// Config holds configuration used when initializing Elastic APM
// Contains connection details, authentication credentials, and service metadata
// Environment variables are set based on these values during initialization
//
// Config 保存 Elastic APM 初始化配置
// 包含连接详情、认证凭据和服务元数据
// 初始化时根据这些值设置环境变量
type Config struct {
	Environment    string   // Environment name (e.g., "production", "staging") // 环境名称（如 "production"、"staging"）
	ServerUrls     []string // Multiple APM servers URLs // 多个 APM 服务地址
	ServerUrl      string   // Single APM URL // 单个 APM 地址
	ApiKey         string   `json:"-"` // API key used in authentication (not logged) // API 密钥用于认证（不记录日志）
	SecretToken    string   `json:"-"` // Secret token used in authentication (not logged) // Secret Token 用于认证（不记录日志）
	ServiceName    string   // Service name used in identification // 服务名称用于标识
	ServiceVersion string   // Service version // 服务版本
	NodeName       string   // Node name used in multi-instance services // 节点名称用于多实例服务
	ServerCertPath string   // Path to certificate used in TLS // TLS 证书路径
	SkipShortSpans bool     // Skip spans when duration is too short // 跳过持续时间过短的 span
}

// SetEnv configures environment variables based on config values
// Uses EnvOption to decide on overriding existing variables
// Maps config fields to ELASTIC_APM_* environment variables
//
// SetEnv 根据配置值设置环境变量
// 使用 EnvOption 决定是否覆盖现有变量
// 将配置字段映射到 ELASTIC_APM_* 环境变量
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
	// If not skipping short spans, set minimum duration to 0ms
	// Note: Some components might drop spans with durations less than 1ms due to precision limits
	// 如果不跳过短 span，设置最小持续时间 0ms
	// 注意：某些组件由于精度限制可能丢弃 1ms 以下的 span
	if !cfg.SkipShortSpans {
		evo.SetEnv("ELASTIC_APM_SPAN_STACK_TRACE_MIN_DURATION", "0ms")
	}
}

// EnvOption controls how environment variables are set during initialization
// When Override is false, existing environment variables are kept unchanged
//
// EnvOption 控制初始化时如何设置环境变量
// 当 Override 是 false 时，保留现有的环境变量
type EnvOption struct {
	Override bool // Override existing env vars when true // 当值是 true 时覆盖现有环境变量
}

// NewEnvOption creates a new EnvOption with Override enabled
// Returns instance that overrides existing environment variables
//
// NewEnvOption 创建启用 Override 的新 EnvOption
// 返回会覆盖现有环境变量的实例
func NewEnvOption() *EnvOption {
	return &EnvOption{
		Override: true,
	}
}

// SetEnv sets an environment variable based on the Override setting
// When Override is true, sets the variable without conditions
// When Override is false, sets the variable if it is not yet defined
//
// SetEnv 根据 Override 设置设定环境变量
// 当 Override 是 true 时，无条件设置变量
// 当 Override 是 false 时，在变量未定义时才设置
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
