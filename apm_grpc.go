package elasticapm

import (
	"context"
	"strings"

	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"google.golang.org/grpc/metadata"
)

// StartApmTraceGrpcOutgoingCtx starts a new APM transaction and injects trace context into gRPC metadata
// Returns the transaction and context prepared to make outgoing gRPC calls
// The context contains W3C trace headers enabling distributed tracing
//
// StartApmTraceGrpcOutgoingCtx 启动新的 APM 事务并将追踪上下文注入 gRPC 元数据
// 返回事务和准备好的上下文用于出站 gRPC 调用
// 上下文包含 W3C 追踪头实现分布式追踪
func StartApmTraceGrpcOutgoingCtx(ctx context.Context, name, apmTxnType string) (*apm.Transaction, context.Context) {
	apmTransaction := apm.DefaultTracer().StartTransaction(name, apmTxnType)
	ctx = ContextWithTraceGrpcOutgoing(ctx, apmTransaction)
	return apmTransaction, ctx
}

// ContextWithTraceGrpcOutgoing attaches transaction to context and adds gRPC outgoing trace headers
// First adds the transaction to context, then injects trace metadata
// Propagates traces across gRPC service boundaries
//
// ContextWithTraceGrpcOutgoing 将事务附加到上下文并添加 gRPC 出站追踪头
// 首先将事务添加到上下文，然后注入追踪元数据
// 跨 gRPC 服务边界传播追踪
func ContextWithTraceGrpcOutgoing(ctx context.Context, apmTransaction *apm.Transaction) context.Context {
	ctx = apm.ContextWithTransaction(ctx, apmTransaction)
	ctx = ContextWithGrpcOutgoingTrace(ctx, apmTransaction.TraceContext())
	return ctx
}

// ContextWithGrpcOutgoingTrace injects W3C trace headers into gRPC outgoing metadata
// Adds traceparent and tracestate headers enabling distributed tracing
// Headers are set in lowercase ensuring correct gRPC metadata handling
//
// ContextWithGrpcOutgoingTrace 将 W3C 追踪头注入 gRPC 出站元数据
// 添加 traceparent 和 tracestate 头实现分布式追踪
// 头信息使用小写确保正确处理 gRPC 元数据
func ContextWithGrpcOutgoingTrace(ctx context.Context, apmTraceContext apm.TraceContext) context.Context {
	upk := apmhttp.FormatTraceparentHeader(apmTraceContext)
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs(strings.ToLower(apmhttp.W3CTraceparentHeader), upk)
	} else {
		md = md.Copy()
		md.Set(strings.ToLower(apmhttp.W3CTraceparentHeader), upk)
	}
	if stateMessage := apmTraceContext.State.String(); stateMessage != "" {
		md.Set(strings.ToLower(apmhttp.TracestateHeader), stateMessage)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// SetLog sets a custom logging implementation with APM tracing
// Replaces the default logging with the provided implementation
//
// SetLog 设置自定义日志实现与 APM 追踪
// 用提供的实现替换默认日志
func SetLog(LOG apm.Logger) {
	apm.DefaultTracer().SetLogger(LOG)
}

// Close flushes pending data and closes the APM tracing
// Should be called before application shutdown to ensure data is sent
// Logs tracing statistics when monitoring is needed
//
// Close 刷新待处理数据并关闭 APM 追踪
// 应在应用关闭前调用以确保数据发送完成
// 在需要监控时记录追踪统计信息
func Close() {
	tracing := apm.DefaultTracer()
	// Flush data to ensure sent before exit
	// 刷新数据确保在退出前发送
	tracing.Flush(nil)
	tracing.Close()
	zaplog.SUG.Debug(neatjsons.S(tracing.Stats()))
}
