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

func StartRequestTransaction(name string) *apm.Transaction {
	return apm.DefaultTracer().StartTransaction(name, "request")
}

func StartTransaction(name, apmTxnType string) (apmTxn *apm.Transaction) {
	return apm.DefaultTracer().StartTransaction(name, apmTxnType)
}

func ContextWithTransaction(ctx context.Context, apmTxn *apm.Transaction) context.Context {
	return apm.ContextWithTransaction(ctx, apmTxn)
}

func TransactionFromContext(ctx context.Context) (transaction *apm.Transaction) {
	return apm.TransactionFromContext(ctx)
}

func StartApmTxnTraceGrpcOutgoing(ctx context.Context, name, apmTxnType string) (*apm.Transaction, context.Context) {
	apmTxn := StartTransaction(name, apmTxnType)
	ctx = ContextWithTraceGrpcOutgoing(ctx, apmTxn)
	return apmTxn, ctx
}

func ContextWithTraceGrpcOutgoing(ctx context.Context, apmTxn *apm.Transaction) context.Context {
	ctx = ContextWithTransaction(ctx, apmTxn)
	ctx = ContextWithGrpcOutgoingTrace(ctx, apmTxn.TraceContext())
	return ctx
}

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

func SetLog(LOG apm.Logger) {
	apm.DefaultTracer().SetLogger(LOG)
}

func Close() {
	tracer := apm.DefaultTracer()
	tracer.Flush(nil) //直接带个刷缓冲即可，否则程序快速退出时没数据
	tracer.Close()
	zaplog.SUG.Debug(neatjsons.S(tracer.Stats()))
}
