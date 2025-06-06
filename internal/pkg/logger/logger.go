package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

// Init 初始化日志
func Init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// 根据环境变量决定日志格式
	if os.Getenv("LOG_FORMAT") == "json" {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		Logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	// 设置为默认日志
	slog.SetDefault(Logger)
}

// Info 信息日志
func Info(msg string, args ...interface{}) {
	Logger.Info(msg, args...)
}

// Error 错误日志
func Error(msg string, args ...interface{}) {
	Logger.Error(msg, args...)
}

// Warn 警告日志
func Warn(msg string, args ...interface{}) {
	Logger.Warn(msg, args...)
}

// Debug 调试日志
func Debug(msg string, args ...interface{}) {
	Logger.Debug(msg, args...)
}
