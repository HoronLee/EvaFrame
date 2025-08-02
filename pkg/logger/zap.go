package logger

import (
	"encoding/json"
	"os"
	"path/filepath"

	"evaframe/pkg/config"

	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ProviderSet = wire.NewSet(NewLogger)

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

// Dump 调试专用，不会中断程序，会在终端打印出 warning 消息。
// 第一个参数会使用 json.Marshal 进行渲染，第二个参数消息（可选）
//
// logger.Dump(user.User{Name:"test"})
// logger.Dump(user.User{Name:"test"}, "用户信息")
func (l *Logger) Dump(value interface{}, msg ...string) {
	valueString := l.jsonString(value)
	// 判断第二个参数是否传参 msg
	if len(msg) > 0 {
		l.Warn("Dump", zap.String(msg[0], valueString))
	} else {
		l.Warn("Dump", zap.String("data", valueString))
	}
}

// LogIf 当 err != nil 时记录 error 等级的日志
func (l *Logger) LogIf(msg string, err error) {
	if err != nil {
		l.Error(msg, zap.Error(err))
	}
}

// LogWarnIf 当 err != nil 时记录 warning 等级的日志
func (l *Logger) LogWarnIf(msg string, err error) {
	if err != nil {
		l.Warn(msg, zap.Error(err))
	}
}

// LogInfoIf 当 err != nil 时记录 info 等级的日志
func (l *Logger) LogInfoIf(msg string, err error) {
	if err != nil {
		l.Info(msg, zap.Error(err))
	}
}

// DebugString 记录一条字符串类型的 debug 日志
func (l *Logger) DebugString(moduleName, name, msg string) {
	l.Debug(moduleName, zap.String(name, msg))
}

// InfoString 记录一条字符串类型的 info 日志
func (l *Logger) InfoString(moduleName, name, msg string) {
	l.Info(moduleName, zap.String(name, msg))
}

// WarnString 记录一条字符串类型的 warn 日志
func (l *Logger) WarnString(moduleName, name, msg string) {
	l.Warn(moduleName, zap.String(name, msg))
}

// ErrorString 记录一条字符串类型的 error 日志
func (l *Logger) ErrorString(moduleName, name, msg string) {
	l.Error(moduleName, zap.String(name, msg))
}

// FatalString 记录一条字符串类型的 fatal 日志
func (l *Logger) FatalString(moduleName, name, msg string) {
	l.Fatal(moduleName, zap.String(name, msg))
}

// DebugJSON 记录对象类型的 debug 日志
func (l *Logger) DebugJSON(moduleName, name string, value interface{}) {
	l.Debug(moduleName, zap.String(name, l.jsonString(value)))
}

// InfoJSON 记录对象类型的 info 日志
func (l *Logger) InfoJSON(moduleName, name string, value interface{}) {
	l.Info(moduleName, zap.String(name, l.jsonString(value)))
}

// WarnJSON 记录对象类型的 warn 日志
func (l *Logger) WarnJSON(moduleName, name string, value interface{}) {
	l.Warn(moduleName, zap.String(name, l.jsonString(value)))
}

// ErrorJSON 记录对象类型的 error 日志
func (l *Logger) ErrorJSON(moduleName, name string, value interface{}) {
	l.Error(moduleName, zap.String(name, l.jsonString(value)))
}

// FatalJSON 记录对象类型的 fatal 日志
func (l *Logger) FatalJSON(moduleName, name string, value interface{}) {
	l.Fatal(moduleName, zap.String(name, l.jsonString(value)))
}

func (l *Logger) jsonString(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		l.Error("Logger", zap.String("JSON marshal error", err.Error()))
	}
	return string(b)
}