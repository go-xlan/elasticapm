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
	// 初始化 zap 日志
	zaplog.SUG.Info("Starting gRPC APM demo")

	// Configure APM
	// 配置 APM
	cfg := &elasticapm.Config{
		Environment:    "development",
		ServerUrl:      "http://localhost:8200",
		ServiceName:    "demo-grpc-client",
		ServiceVersion: "1.0.0",
	}

	// Initialize APM
	// 初始化 APM
	must.Done(elasticapm.Initialize(cfg))
	defer elasticapm.Close()

	zaplog.SUG.Info("APM initialized for gRPC demo")

	// Simulate gRPC client call with distributed tracing
	// 模拟带分布式追踪的 gRPC 客户端调用
	ctx := context.Background()
	callRemoteService(ctx)

	zaplog.SUG.Info("gRPC demo completed")
}

func callRemoteService(ctx context.Context) {
	// Start APM transaction with gRPC outgoing context
	// The trace context will be auto injected into gRPC metadata
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

	// Simulate preparing request
	// 模拟准备请求
	prepareSpan := txn.StartSpan("prepare-request", "internal", nil)
	prepareRequest()
	prepareSpan.End()

	// Simulate gRPC call
	// In production code, you would pass tracedCtx to your gRPC client call:
	// response, err := grpcClient.Method(tracedCtx, request)
	// 模拟 gRPC 调用
	// 在生产代码中，你会将 tracedCtx 传递给 gRPC 客户端调用：
	// response, err := grpcClient.Method(tracedCtx, request)
	grpcSpan := txn.StartSpan("grpc.Call", "external.grpc", nil)
	simulateGrpcCall(tracedCtx)
	grpcSpan.End()

	// Simulate processing response
	// 模拟处理响应
	processSpan := txn.StartSpan("process-response", "internal", nil)
	processResponse()
	processSpan.End()

	zaplog.SUG.Info("gRPC call completed")
}

func prepareRequest() {
	// Simulate request preparation
	// 模拟请求准备
	time.Sleep(20 * time.Millisecond)
	zaplog.SUG.Debug("Request prepared")
}

func simulateGrpcCall(ctx context.Context) {
	// Simulate gRPC network call
	// The trace context is already in the metadata thanks to StartApmTraceGrpcOutgoingCtx
	// 模拟 gRPC 网络调用
	// 追踪上下文已通过 StartApmTraceGrpcOutgoingCtx 放入元数据中
	time.Sleep(100 * time.Millisecond)
	zaplog.SUG.Debug("gRPC call executed")
}

func processResponse() {
	// Simulate response processing
	// 模拟响应处理
	time.Sleep(30 * time.Millisecond)
	zaplog.SUG.Debug("Response processed")
}
