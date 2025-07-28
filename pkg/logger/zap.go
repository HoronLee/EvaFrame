package logger

import (
	"os"
	"path/filepath"

	"evaframe/pkg/config"

	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	// 确保日志目录存在
	logDir := filepath.Dir(cfg.Logger.LogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// 解析日志级别
	var level zapcore.Level
	switch cfg.Logger.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置 Zap
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: cfg.Server.Mode == "debug",
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", cfg.Logger.LogPath},
		ErrorOutputPaths: []string{"stderr"},
	}

	return config.Build()
}

var ProviderSet = wire.NewSet(NewLogger)
