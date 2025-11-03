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
	// 首先初始化 zap 日志
	zaplog.SUG.Info("Starting APM demo")

	// Configure APM settings
	// 配置 APM 设置
	cfg := &elasticapm.Config{
		Environment:    "development",
		ServerUrl:      "http://localhost:8200", // Your APM server URL // 你的 APM 服务器地址
		ServiceName:    "demo-basic-service",
		ServiceVersion: "1.0.0",
		SkipShortSpans: false, // Capture all spans // 捕获所有 span
	}

	// Initialize APM
	// 初始化 APM
	must.Done(elasticapm.Initialize(cfg))
	defer elasticapm.Close()

	zaplog.SUG.Info("APM initialized", zap.String("version", elasticapm.GetApmAgentVersion()))

	// Verify version compatibility
	// 验证版本兼容性
	if elasticapm.CheckApmAgentVersion(apm.AgentVersion) {
		zaplog.SUG.Info("APM version check passed")
	}

	// Start a transaction
	// 启动一个事务
	txn := apm.DefaultTracer().StartTransaction("demo-operation", "custom")
	defer txn.End()

	zaplog.SUG.Info("Transaction started", zap.String("transaction_id", txn.TraceContext().Trace.String()))

	// Simulate some work with a span
	// 使用 span 模拟一些工作
	span := txn.StartSpan("process-data", "internal", nil)
	processData()
	span.End()

	// Second span simulating database operation
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
	// Simulate data processing
	// 模拟数据处理
	time.Sleep(100 * time.Millisecond)
	zaplog.SUG.Debug("Data processed")
}

func simulateDatabaseOperation() {
	// Simulate database operation
	// 模拟数据库操作
	time.Sleep(50 * time.Millisecond)
	zaplog.SUG.Debug("Database operation executed")
}
