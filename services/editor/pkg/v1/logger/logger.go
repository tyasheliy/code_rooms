package logger

import (
	"context"
	"io"
	"log/slog"
)

const (
	TEXT_FORMAT = "text"
	JSON_FORMAT = "json"

	DEBUG_LEVEL = "debug"
	INFO_LEVEL  = "info"
)

type AppLogger interface {
	Debug(ctx context.Context, message string, args ...any)
	Info(ctx context.Context, message string, args ...any)
	Warn(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, message string, args ...any)
}

type AppLoggerOptions struct {
	Format string
	Level  string
	Output io.Writer
}

type SlogAppLogger struct {
	l *slog.Logger
}

func NewSlogAppLogger(options AppLoggerOptions) *SlogAppLogger {
	var lvl slog.Level
	switch options.Level {
	case DEBUG_LEVEL:
		lvl = slog.LevelDebug
		break
	case INFO_LEVEL:
		lvl = slog.LevelInfo
	default:
		lvl = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		Level: lvl,
	}

	var handler slog.Handler
	switch options.Format {
	case TEXT_FORMAT:
		handler = slog.NewTextHandler(options.Output, handlerOptions)
		break
	case JSON_FORMAT:
		handler = slog.NewJSONHandler(options.Output, handlerOptions)
	}

	return &SlogAppLogger{
		l: slog.New(handler),
	}
}

func (l *SlogAppLogger) Debug(ctx context.Context, message string, args ...any) {
	l.l.DebugContext(ctx, message, args...)
}

func (l *SlogAppLogger) Info(ctx context.Context, message string, args ...any) {
	l.l.InfoContext(ctx, message, args...)
}

func (l *SlogAppLogger) Warn(ctx context.Context, message string, args ...any) {
	l.l.WarnContext(ctx, message, args...)
}

func (l *SlogAppLogger) Error(ctx context.Context, message string, args ...any) {
	l.l.ErrorContext(ctx, message, args...)
}
