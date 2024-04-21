package logging

import (
	"context"
	"os"
	"runtime"
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
	logger.Log(msg, 6, mergeMaps(collectData(), mergeMaps(data...)))
}

func Error(ctx context.Context, msg string, data ...map[string]interface{}) {
	logger.Log(msg, 3, mergeMaps(collectData(), mergeMaps(data...)))
}

func Fatal(ctx context.Context, msg string, data ...map[string]interface{}) {
	logger.Log(msg, 1, mergeMaps(collectData(), mergeMaps(data...)))
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

func collectData() map[string]interface{} {
	_, file, no, _ := runtime.Caller(2)
	return map[string]interface{}{
		"FILE": file,
		"LINE": no,
	}
}
