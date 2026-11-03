// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

// Package logger 原来是 zap + lumberjack，裁剪后改成标准库 log/slog：
// 日志写到 slog.Default()，落点、格式、滚动都由宿主程序配。函数名保持原样，调用方不用改。
package logger

import (
	stdctx "context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"goadmin/context"
	"goadmin/modules/trace"
)

var logger = &Logger{Level: slog.LevelInfo}

type Logger struct {
	infoLogOff   bool
	errorLogOff  bool
	accessLogOff bool

	accessAssetsLogOff bool

	sqlLogOpen bool

	Level slog.Level
}

type Config struct {
	InfoLogOff   bool
	ErrorLogOff  bool
	AccessLogOff bool

	SqlLogOpen bool

	AccessAssetsLogOff bool

	// Level -1 debug、0 info、1 warn、2 error（和原来的配置一样）
	Level int8
}

func InitWithConfig(cfg Config) {
	logger.infoLogOff = cfg.InfoLogOff
	logger.errorLogOff = cfg.ErrorLogOff
	logger.accessLogOff = cfg.AccessLogOff
	logger.sqlLogOpen = cfg.SqlLogOpen
	logger.accessAssetsLogOff = cfg.AccessAssetsLogOff
	logger.Level = levelOf(cfg.Level)
}

func levelOf(level int8) slog.Level {
	switch level {
	case -1:
		return slog.LevelDebug
	case 1:
		return slog.LevelWarn
	case 2:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// OpenSQLLog set the sqlLogOpen true.
func OpenSQLLog() {
	logger.sqlLogOpen = true
}

func log(level slog.Level, msg string, args ...any) {
	if level < logger.Level {
		return
	}
	if level >= slog.LevelError && logger.errorLogOff || level < slog.LevelError && logger.infoLogOff {
		return
	}
	slog.Default().Log(stdctx.Background(), level, msg, append([]any{"logger", "goadmin"}, args...)...)
}

func logCtx(ctx *context.Context, level slog.Level, format string, args ...interface{}) {
	log(level, fmt.Sprintf(format, args...), "traceID", trace.GetTraceID(ctx))
}

// Debug print the debug message.
func Debug(info ...interface{}) { log(slog.LevelDebug, fmt.Sprint(info...)) }

// Debugf print the debug message.
func Debugf(template string, args ...interface{}) {
	log(slog.LevelDebug, fmt.Sprintf(template, args...))
}

// Info print the info message.
func Info(info ...interface{}) { log(slog.LevelInfo, fmt.Sprint(info...)) }

// InfoCtx print the info message with ctx.
func InfoCtx(ctx *context.Context, format string, args ...interface{}) {
	logCtx(ctx, slog.LevelInfo, format, args...)
}

// Infof print the info message.
func Infof(template string, args ...interface{}) {
	log(slog.LevelInfo, fmt.Sprintf(template, args...))
}

// Warn print the warning message.
func Warn(info ...interface{}) { log(slog.LevelWarn, fmt.Sprint(info...)) }

// WarnCtx print the warning message with ctx.
func WarnCtx(ctx *context.Context, format string, args ...interface{}) {
	logCtx(ctx, slog.LevelWarn, format, args...)
}

// Warnf print the warning message.
func Warnf(template string, args ...interface{}) {
	log(slog.LevelWarn, fmt.Sprintf(template, args...))
}

// Error print the error message.
func Error(err ...interface{}) { log(slog.LevelError, fmt.Sprint(err...)) }

// ErrorCtx print the error message with ctx.
func ErrorCtx(ctx *context.Context, format string, args ...interface{}) {
	logCtx(ctx, slog.LevelError, format, args...)
}

// Errorf print the error message.
func Errorf(template string, args ...interface{}) {
	log(slog.LevelError, fmt.Sprintf(template, args...))
}

// Fatal print the fatal message and exit.
func Fatal(info ...interface{}) {
	log(slog.LevelError, fmt.Sprint(info...))
	os.Exit(1)
}

// Fatalf print the fatal message and exit.
func Fatalf(template string, args ...interface{}) {
	log(slog.LevelError, fmt.Sprintf(template, args...))
	os.Exit(1)
}

// Panic print the panic message and panic.
func Panic(info ...interface{}) {
	msg := fmt.Sprint(info...)
	log(slog.LevelError, msg)
	panic(msg)
}

// Panicf print the panic message and panic.
func Panicf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	log(slog.LevelError, msg)
	panic(msg)
}

// Access print the access message.
func Access(ctx *context.Context) {
	if logger.accessLogOff || logger.accessAssetsLogOff && filepath.Ext(ctx.Path()) != "" {
		return
	}
	log(slog.LevelInfo, "access", "traceID", trace.GetTraceID(ctx), "status", ctx.Response.StatusCode,
		"method", ctx.Method(), "path", ctx.Path())
}

// LogSQL print the sql info message.
func LogSQL(statement string, args []interface{}) {
	if logger.sqlLogOpen && statement != "" {
		log(slog.LevelInfo, "sql", "statement", statement, "args", args)
	}
}
