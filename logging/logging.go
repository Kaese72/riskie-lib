package logging

import (
	"context"
	"os"
	"runtime"

	"go.elastic.co/apm"
)

type Logger interface {
	Log(string, int, ...map[string]interface{})
}

var logger Logger = JSONLogger{}

func SetLogger(newLogger Logger) {
	logger = newLogger
}

var debugLogging bool = false

func Debug(ctx context.Context, msg string, data ...map[string]interface{}) {
	if debugLogging {
		logger.Log(msg, 7, data...)
	}
}

func Info(ctx context.Context, msg string, data ...map[string]interface{}) {
	logger.Log(msg, 6, mergeMaps(collectData(ctx), mergeMaps(data...)))
}

func Error(ctx context.Context, msg string, data ...map[string]interface{}) {
	logger.Log(msg, 3, mergeMaps(collectData(ctx), mergeMaps(data...)))
}

func Fatal(ctx context.Context, msg string, data ...map[string]interface{}) {
	logger.Log(msg, 1, mergeMaps(collectData(ctx), mergeMaps(data...)))
	os.Exit(1)
}

func SetDebugLogging(flag bool) {
	debugLogging = flag
}

func mergeMaps(datas ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, data := range datas {
		for key, value := range data {
			merged[key] = value
		}
	}
	return merged
}

func collectData(ctx context.Context) map[string]interface{} {
	labels := map[string]interface{}{}
	// Completely stolen from documentation, https://www.elastic.co/guide/en/apm/agent/go/current/log-correlation-ids.html
	// Some slight modifications to create correct types
	tx := apm.TransactionFromContext(ctx)
	if tx != nil {
		traceContext := tx.TraceContext()
		labels["trace.id"] = traceContext.Trace.String()
		labels["transaction.id"] = traceContext.Span.String()
		if span := apm.SpanFromContext(ctx); span != nil {
			labels["span.id"] = span.TraceContext().Span.String()
		}
	}
	_, file, no, _ := runtime.Caller(2)
	labels["FILE"] = file
	labels["LINE"] = no
	return labels
}
