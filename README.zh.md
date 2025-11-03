[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-xlan/elasticapm/release.yml?branch=main&label=BUILD)](https://github.com/go-xlan/elasticapm/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-xlan/elasticapm)](https://pkg.go.dev/github.com/go-xlan/elasticapm)
[![Coverage Status](https://img.shields.io/coveralls/github/go-xlan/elasticapm/main.svg)](https://coveralls.io/github/go-xlan/elasticapm?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.23+-lightgrey.svg)](https://github.com/go-xlan/elasticapm)
[![GitHub Release](https://img.shields.io/github/release/go-xlan/elasticapm.svg)](https://github.com/go-xlan/elasticapm/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-xlan/elasticapm)](https://goreportcard.com/report/github.com/go-xlan/elasticapm)

# elasticapm

基于 `go.elastic.co/apm/v2` 的简洁优雅的 Elastic APM（应用性能监控）Go 封装库。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🎯 **简洁 APM 配置**: 基于配置结构体的初始化，减少样板代码
⚡ **Zap 日志集成**: 内置 Zap 日志支持，包含 APM 上下文追踪
🔄 **gRPC 分布式追踪**: W3C 追踪标头在 gRPC 边界间的传播
🌍 **环境变量**: 自动环境设置，支持覆盖控制
📋 **版本匹配**: 确保使用 v2，防止 v1/v2 混用陷阱

## 安装

```bash
go get github.com/go-xlan/elasticapm
```

### 依赖要求

- Go 1.23.0 及以上版本
- Elastic APM Server v2.x

## 使用示例

### 基础 APM 与事务和跨度

此示例展示完整的 APM 设置，包含事务追踪和跨度埋点：

```go
package main

import (
	"time"

	"github.com/go-xlan/elasticapm"
	"github.com/yyle88/must"
	"github.com/yyle88/zaplog"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

func main() {
	// 首先初始化 zap 日志
	zaplog.SUG.Info("Starting APM demo")

	// 配置 APM 设置
	cfg := &elasticapm.Config{
		Environment:    "development",
		ServerUrl:      "http://localhost:8200",
		ServiceName:    "demo-basic-service",
		ServiceVersion: "1.0.0",
		SkipShortSpans: false, // 捕获所有 span
	}

	// 初始化 APM
	must.Done(elasticapm.Initialize(cfg))
	defer elasticapm.Close()

	zaplog.SUG.Info("APM initialized", zap.String("version", elasticapm.GetApmAgentVersion()))

	// 验证版本兼容性
	if elasticapm.CheckApmAgentVersion(apm.AgentVersion) {
		zaplog.SUG.Info("APM version check passed")
	}

	// 启动一个事务
	txn := apm.DefaultTracer().StartTransaction("demo-operation", "custom")
	defer txn.End()

	zaplog.SUG.Info("Transaction started", zap.String("transaction_id", txn.TraceContext().Trace.String()))

	// 使用 span 模拟一些工作
	span := txn.StartSpan("process-data", "internal", nil)
	processData()
	span.End()

	// 第二个 span 用于模拟数据库操作
	dbSpan := txn.StartSpan("database-query", "db.query", nil)
	dbSpan.Context.SetDatabase(apm.DatabaseSpanContext{
		Statement: "SELECT * FROM users WHERE id = ?",
		Type:      "sql",
	})
	simulateDatabaseOperation()
	dbSpan.End()

	zaplog.SUG.Info("Demo completed")
}

func processData() {
	// 模拟数据处理
	time.Sleep(100 * time.Millisecond)
	zaplog.SUG.Debug("Data processed")
}

func simulateDatabaseOperation() {
	// 模拟数据库操作
	time.Sleep(50 * time.Millisecond)
	zaplog.SUG.Debug("Database operation executed")
}
```

⬆️ **源码：** [源码](internal/demos/demo1x/main.go)

### gRPC 分布式追踪

此示例演示 W3C 追踪标头在 gRPC 服务边界间的传播：

```go
package main

import (
	"context"
	"time"

	"github.com/go-xlan/elasticapm"
	"github.com/yyle88/must"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func main() {
	// 初始化 zap 日志
	zaplog.SUG.Info("Starting gRPC APM demo")

	// 配置 APM
	cfg := &elasticapm.Config{
		Environment:    "development",
		ServerUrl:      "http://localhost:8200",
		ServiceName:    "demo-grpc-client",
		ServiceVersion: "1.0.0",
	}

	// 初始化 APM
	must.Done(elasticapm.Initialize(cfg))
	defer elasticapm.Close()

	zaplog.SUG.Info("APM initialized for gRPC demo")

	// 模拟带分布式追踪的 gRPC 客户端调用
	ctx := context.Background()
	callRemoteService(ctx)

	zaplog.SUG.Info("gRPC demo completed")
}

func callRemoteService(ctx context.Context) {
	// 启动带 gRPC 出站上下文的 APM 事务
	// 追踪上下文将自动注入到 gRPC 元数据中
	txn, tracedCtx := elasticapm.StartApmTraceGrpcOutgoingCtx(
		ctx,
		"grpc-call-remote-service",
		"request",
	)
	defer txn.End()

	zaplog.SUG.Info("Starting gRPC call",
		zap.String("trace_id", txn.TraceContext().Trace.String()),
		zap.String("transaction_id", txn.TraceContext().Span.String()),
	)

	// 模拟准备请求
	prepareSpan := txn.StartSpan("prepare-request", "internal", nil)
	prepareRequest()
	prepareSpan.End()

	// 模拟 gRPC 调用
	// 在生产代码中，你会将 tracedCtx 传递给 gRPC 客户端调用：
	// response, err := grpcClient.Method(tracedCtx, request)
	grpcSpan := txn.StartSpan("grpc.Call", "external.grpc", nil)
	simulateGrpcCall(tracedCtx)
	grpcSpan.End()

	// 模拟处理响应
	processSpan := txn.StartSpan("process-response", "internal", nil)
	processResponse()
	processSpan.End()

	zaplog.SUG.Info("gRPC call completed")
}

func prepareRequest() {
	// 模拟请求准备
	time.Sleep(20 * time.Millisecond)
	zaplog.SUG.Debug("Request prepared")
}

func simulateGrpcCall(ctx context.Context) {
	// 模拟 gRPC 网络调用
	// 追踪上下文已通过 StartApmTraceGrpcOutgoingCtx 放入元数据中
	time.Sleep(100 * time.Millisecond)
	zaplog.SUG.Debug("gRPC call executed")
}

func processResponse() {
	// 模拟响应处理
	time.Sleep(30 * time.Millisecond)
	zaplog.SUG.Debug("Response processed")
}
```

⬆️ **源码：** [源码](internal/demos/demo2x/main.go)

## 配置选项

| 字段 | 类型 | 说明 |
|------|------|------|
| `Environment` | `string` | 环境名称（如 "production"、"staging"） |
| `ServerUrl` | `string` | 单个 APM 服务器地址 |
| `ServerUrls` | `[]string` | 多个 APM 服务器地址 |
| `ApiKey` | `string` | 用于 APM 服务器认证的 API 密钥 |
| `SecretToken` | `string` | 用于 APM 服务器认证的密钥令牌 |
| `ServiceName` | `string` | 标识此服务的名称 |
| `ServiceVersion` | `string` | 服务版本 |
| `NodeName` | `string` | 多实例服务的节点名称 |
| `ServerCertPath` | `string` | 服务器证书路径 |
| `SkipShortSpans` | `bool` | 跳过短于阈值的 span |

## API 参考

### 核心函数

- `Initialize(cfg *Config) error` - 使用默认选项初始化 APM
- `InitializeWithOptions(cfg *Config, evo *EnvOption, setEnvs ...func()) error` - 使用自定义选项初始化
- `Close()` - 刷新并关闭 APM 追踪
- `SetLog(LOG apm.Logger)` - 设置自定义日志

### 版本函数

- `GetApmAgentVersion() string` - 获取当前 APM agent 版本
- `CheckApmAgentVersion(agentVersion string) bool` - 验证版本匹配

### gRPC 函数

- `StartApmTraceGrpcOutgoingCtx(ctx, name, apmTxnType) (*apm.Transaction, context.Context)` - 启动带追踪的 gRPC 调用
- `ContextWithTraceGrpcOutgoing(ctx, apmTransaction) context.Context` - 向上下文添加追踪
- `ContextWithGrpcOutgoingTrace(ctx, apmTraceContext) context.Context` - 向出站元数据添加追踪上下文

## 高级用法

### 环境变量配置

包支持通过环境变量进行配置。你可以控制是否覆盖已有变量：

```go
cfg := &elasticapm.Config{
    Environment:    "production",
    ServerUrl:      "http://localhost:8200",
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
}

envOption := &elasticapm.EnvOption{
    Override: true, // 覆盖已有的环境变量
}

must.Done(elasticapm.InitializeWithOptions(cfg, envOption))
defer elasticapm.Close()
```

### 自定义日志集成

与自定义 zap 日志设置集成：

```go
import (
    "github.com/go-xlan/elasticapm"
    "github.com/go-xlan/elasticapm/apmzaplog"
)

// 使用自定义 zap 日志初始化 APM
must.Done(elasticapm.Initialize(cfg))
elasticapm.SetLog(apmzaplog.NewLog())
defer elasticapm.Close()
```

### 上下文传播

在微服务架构中，需要在服务间传递追踪上下文：

```go
// 服务 A: 启动事务
txn := apm.DefaultTracer().StartTransaction("external-call", "request")
ctx := apm.ContextWithTransaction(context.Background(), txn)

// 将追踪上下文注入 gRPC 元数据
ctx = elasticapm.ContextWithTraceGrpcOutgoing(ctx, txn)

// 发起 gRPC 调用
response := grpcClient.Method(ctx, request)

txn.End()
```

## 最佳实践

### 服务命名

选择能够反映实际功能的服务名称：

- 使用小写字母和连字符：`user-service`、`payment-gateway`
- 当多个团队共享基础设施时包含团队名称：`team-a-user-service`
- 保持名称简洁但具有描述性

### 环境配置

设置不同的环境以分离追踪数据：

- `development` - 在你的机器上开发
- `staging` - 预生产测试
- `production` - 生产环境流量
- `testing` - 自动化测试运行

### 性能优化

在高吞吐量应用中减少开销：

```go
cfg := &elasticapm.Config{
    Environment:    "production",
    ServerUrl:      "http://localhost:8200",
    ServiceName:    "my-service",
    SkipShortSpans: true, // 跳过短于阈值的 span
}
```

### 版本追踪

始终在服务配置中包含语义版本：

```go
cfg := &elasticapm.Config{
    ServiceVersion: "1.2.3", // 语义版本
}
```

这有助于将性能变化与部署关联起来。

## 故障排查

### 连接问题

如果 APM 无法连接到服务器：

1. 验证服务器 URL 是否可访问
2. 检查防火墙设置是否允许出站连接
3. 验证 API 密钥是否正确
4. 检查 APM 服务器日志

### 缺失追踪数据

如果追踪数据没有出现在 Kibana 中：

1. 确保事务已结束：`defer txn.End()`
2. 验证服务名称与 Kibana 过滤器匹配
3. 检查环境设置是否符合预期
4. 在退出前调用 `Close()` 刷新待处理数据

### 版本不匹配

如果看到关于缺失符号的错误：

```bash
# 检查所有依赖是否使用 v2
go list -m all | grep elastic
```

确保所有导入使用 `go.elastic.co/apm/v2`，而非 `go.elastic.co/apm`。

## 版本兼容性

本包**需要 v2**：`go.elastic.co/apm/v2`，不兼容 v1.x 包 `go.elastic.co/apm`。

版本检查函数确保所有依赖都使用 v2，避免混用 v1 和 v2 包导致的常见问题，这两个包维护着独立的追踪实例。

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-xlan/elasticapm.svg?variant=adaptive)](https://starchart.cc/go-xlan/elasticapm)
