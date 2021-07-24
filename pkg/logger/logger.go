package logger

import (
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log = func() *zap.Logger {
	_, file, _, _ := runtime.Caller(1)
	slash := strings.LastIndex(file, "/")
	file = file[slash+1:]

	var logger *zap.Logger
	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ = config.Build()
	return logger
}()
