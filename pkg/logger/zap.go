package logger

import (
	"os"
	"path/filepath"

	"evaframe/pkg/config"

	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ProviderSet = wire.NewSet(NewLogger)

var _gLogger *Logger

// Init 初始化全局单例日志
func Init(cfg *config.Config) error {
	var err error
	_gLogger, err = NewLogger(cfg)
	if err != nil {
		return err
	}
	return nil
}

// L 返回全局单例日志实例
func L() *Logger {
	return _gLogger
}

// Logger wraps zap.Logger to provide helper functions.
type Logger struct {
	*zap.Logger
}

func NewLogger(cfg *config.Config) (*Logger, error) {
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
	zapCfg := zap.Config{
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

	zlog, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zlog}, nil
}
