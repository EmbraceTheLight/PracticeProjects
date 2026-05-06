package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var atom = zap.NewAtomicLevelAt(zap.DebugLevel)

func init() {
	config := zap.Config{
		Level:       atom,
		Development: false,
		Encoding:    "json", // 使用 json 格式编码
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "caller",
			EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	tempLogger, err := config.Build()
	if err != nil {
		panic(err)
	}
	logger = tempLogger.WithOptions(zap.AddCallerSkip(1))
}

// SetLevel 动态设置日志级别
func SetLevel(level string) {
	tLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		fmt.Printf("Invalid level. Input: %s\n", level)
		return
	}
	atom.SetLevel(tLevel)
}
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
