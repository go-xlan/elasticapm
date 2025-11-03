[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-xlan/elasticapm/release.yml?branch=main&label=BUILD)](https://github.com/go-xlan/elasticapm/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-xlan/elasticapm)](https://pkg.go.dev/github.com/go-xlan/elasticapm)
[![Coverage Status](https://img.shields.io/coveralls/github/go-xlan/elasticapm/main.svg)](https://coveralls.io/github/go-xlan/elasticapm?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.23+-lightgrey.svg)](https://github.com/go-xlan/elasticapm)
[![GitHub Release](https://img.shields.io/github/release/go-xlan/elasticapm.svg)](https://github.com/go-xlan/elasticapm/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-xlan/elasticapm)](https://goreportcard.com/report/github.com/go-xlan/elasticapm)

# elasticapm

Simple and elegant wrapping of Elastic APM (Application Performance Monitoring) in Go, based on `go.elastic.co/apm/v2`.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🎯 **Simple APM Setup**: Config struct-based initialization with less boilerplate code
⚡ **Zap Logging Integration**: Built-in Zap logging support with APM context tracking
🔄 **gRPC Distributed Tracing**: W3C trace headers propagation across gRPC boundaries
🌍 **Environment Variables**: Automatic environment setup with override controls
📋 **Version Matching**: Ensures v2 usage preventing v1/v2 mixing pitfalls

## Installation

```bash
go get github.com/go-xlan/elasticapm
```

### Requirements

- Go 1.23.0 and above
- Elastic APM Server v2.x

## Usage

### Basic APM with Transactions and Spans

This example shows complete APM setup with transaction tracking and span instrumentation:

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
	// Initialize zap logging first
	zaplog.SUG.Info("Starting APM demo")

	// Configure APM settings
	cfg := &elasticapm.Config{
		Environment:    "development",
		ServerUrl:      "http://localhost:8200",
		ServiceName:    "demo-basic-service",
		ServiceVersion: "1.0.0",
		SkipShortSpans: false, // Capture all spans
	}

	// Initialize APM
	must.Done(elasticapm.Initialize(cfg))
	defer elasticapm.Close()

	zaplog.SUG.Info("APM initialized", zap.String("version", elasticapm.GetApmAgentVersion()))

	// Verify version compatibility
	if elasticapm.CheckApmAgentVersion(apm.AgentVersion) {
		zaplog.SUG.Info("APM version check passed")
	}

	// Start a transaction
	txn := apm.DefaultTracer().StartTransaction("demo-operation", "custom")
	defer txn.End()

	zaplog.SUG.Info("Transaction started", zap.String("transaction_id", txn.TraceContext().Trace.String()))

	// Simulate some work with a span
	span := txn.StartSpan("process-data", "internal", nil)
	processData()
	span.End()

	// Second span simulating database operation
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
	// Simulate data processing
	time.Sleep(100 * time.Millisecond)
	zaplog.SUG.Debug("Data processed")
}

func simulateDatabaseOperation() {
	// Simulate database operation
	time.Sleep(50 * time.Millisecond)
	zaplog.SUG.Debug("Database operation executed")
}
```

⬆️ **Source:** [Source](internal/demos/demo1x/main.go)

### gRPC Distributed Tracing

This example demonstrates W3C trace headers propagation across gRPC service boundaries:

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
	// Initialize zap logging
	zaplog.SUG.Info("Starting gRPC APM demo")

	// Configure APM
	cfg := &elasticapm.Config{
		Environment:    "development",
		ServerUrl:      "http://localhost:8200",
		ServiceName:    "demo-grpc-client",
		ServiceVersion: "1.0.0",
	}

	// Initialize APM
	must.Done(elasticapm.Initialize(cfg))
	defer elasticapm.Close()

	zaplog.SUG.Info("APM initialized for gRPC demo")

	// Simulate gRPC client call with distributed tracing
	ctx := context.Background()
	callRemoteService(ctx)

	zaplog.SUG.Info("gRPC demo completed")
}

func callRemoteService(ctx context.Context) {
	// Start APM transaction with gRPC outgoing context
	// The trace context will be auto injected into gRPC metadata
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

	// Simulate preparing request
	prepareSpan := txn.StartSpan("prepare-request", "internal", nil)
	prepareRequest()
	prepareSpan.End()

	// Simulate gRPC call
	// In production code, you would pass tracedCtx to your gRPC client call:
	// response, err := grpcClient.Method(tracedCtx, request)
	grpcSpan := txn.StartSpan("grpc.Call", "external.grpc", nil)
	simulateGrpcCall(tracedCtx)
	grpcSpan.End()

	// Simulate processing response
	processSpan := txn.StartSpan("process-response", "internal", nil)
	processResponse()
	processSpan.End()

	zaplog.SUG.Info("gRPC call completed")
}

func prepareRequest() {
	// Simulate request preparation
	time.Sleep(20 * time.Millisecond)
	zaplog.SUG.Debug("Request prepared")
}

func simulateGrpcCall(ctx context.Context) {
	// Simulate gRPC network call
	// The trace context is already in the metadata thanks to StartApmTraceGrpcOutgoingCtx
	time.Sleep(100 * time.Millisecond)
	zaplog.SUG.Debug("gRPC call executed")
}

func processResponse() {
	// Simulate response processing
	time.Sleep(30 * time.Millisecond)
	zaplog.SUG.Debug("Response processed")
}
```

⬆️ **Source:** [Source](internal/demos/demo2x/main.go)

## Configuration Options

| Field | Type | Description |
|-------|------|-------------|
| `Environment` | `string` | Environment name (e.g., "production", "staging") |
| `ServerUrl` | `string` | Single APM server URL |
| `ServerUrls` | `[]string` | Multiple APM server URLs |
| `ApiKey` | `string` | API key to authenticate with APM server |
| `SecretToken` | `string` | Secret token to authenticate with APM server |
| `ServiceName` | `string` | Name that identifies this service |
| `ServiceVersion` | `string` | Version of this service |
| `NodeName` | `string` | Name of node in multi-instance setup |
| `ServerCertPath` | `string` | Path to server certificate |
| `SkipShortSpans` | `bool` | Skip spans shorter than threshold |

## API Reference

### Core Functions

- `Initialize(cfg *Config) error` - Initialize APM with default options
- `InitializeWithOptions(cfg *Config, evo *EnvOption, setEnvs ...func()) error` - Initialize with custom options
- `Close()` - Flush and close APM tracing
- `SetLog(LOG apm.Logger)` - Set custom logging

### Version Functions

- `GetApmAgentVersion() string` - Get current APM agent version
- `CheckApmAgentVersion(agentVersion string) bool` - Verify version matching

### gRPC Functions

- `StartApmTraceGrpcOutgoingCtx(ctx, name, apmTxnType) (*apm.Transaction, context.Context)` - Start traced gRPC call
- `ContextWithTraceGrpcOutgoing(ctx, apmTransaction) context.Context` - Add trace to context
- `ContextWithGrpcOutgoingTrace(ctx, apmTraceContext) context.Context` - Add trace context to outgoing metadata

## Advanced Usage

### Environment Variable Configuration

The package supports configuration through environment variables. You can control whether to override existing variables:

```go
cfg := &elasticapm.Config{
    Environment:    "production",
    ServerUrl:      "http://localhost:8200",
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
}

envOption := &elasticapm.EnvOption{
    Override: true, // Override existing environment variables
}

must.Done(elasticapm.InitializeWithOptions(cfg, envOption))
defer elasticapm.Close()
```

### Custom Logging Integration

Integrate with custom zap logging setup:

```go
import (
    "github.com/go-xlan/elasticapm"
    "github.com/go-xlan/elasticapm/apmzaplog"
)

// Initialize APM with custom zap logging
must.Done(elasticapm.Initialize(cfg))
elasticapm.SetLog(apmzaplog.NewLog())
defer elasticapm.Close()
```

### Context Propagation

When working with microservices, trace context needs to be passed between services:

```go
// Service A: Start transaction
txn := apm.DefaultTracer().StartTransaction("external-call", "request")
ctx := apm.ContextWithTransaction(context.Background(), txn)

// Inject trace context into gRPC metadata
ctx = elasticapm.ContextWithTraceGrpcOutgoing(ctx, txn)

// Make gRPC call
response := grpcClient.Method(ctx, request)

txn.End()
```

## Best Practices

### Service Naming

Choose service names that reflect their function:

- Use lowercase with hyphens: `user-service`, `payment-gateway`
- Include team name when multiple teams share infrastructure: `team-a-user-service`
- Keep names concise but descriptive

### Environment Configuration

Set distinct environments to separate traces:

- `development` - Development on your machine
- `staging` - Pre-production testing
- `production` - Live production traffic
- `testing` - Automated test runs

### Performance Optimization

To reduce overhead in high-throughput applications:

```go
cfg := &elasticapm.Config{
    Environment:    "production",
    ServerUrl:      "http://localhost:8200",
    ServiceName:    "my-service",
    SkipShortSpans: true, // Skip spans shorter than threshold
}
```

### Version Tracking

Always include semantic version in service configuration:

```go
cfg := &elasticapm.Config{
    ServiceVersion: "1.2.3", // Semantic version
}
```

This helps correlate performance changes with deployments.

## Troubleshooting

### Connection Issues

If APM cannot connect to the server:

1. Verify the server URL is accessible
2. Check firewall settings allow outbound connections
3. Verify API key is correct
4. Check APM server logs

### Missing Traces

If traces are not appearing in Kibana:

1. Ensure transaction is ended: `defer txn.End()`
2. Verify service name matches in Kibana filters
3. Check environment setting matches what you expect
4. Call `Close()` to flush pending data before exit

### Version Mismatch

If you see errors about missing symbols:

```bash
# Check all dependencies use v2
go list -m all | grep elastic
```

Ensure all imports use `go.elastic.co/apm/v2`, not `go.elastic.co/apm`.

## Version Compatibility

This package **requires v2**: `go.elastic.co/apm/v2`. It will not work with the v1.x package `go.elastic.co/apm`.

The version check functions ensure all dependencies use v2 to avoid the common pitfall of mixing v1 and v2 packages, which maintain separate tracing instances.

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-xlan/elasticapm.svg?variant=adaptive)](https://starchart.cc/go-xlan/elasticapm)
