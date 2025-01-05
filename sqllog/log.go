package sqllog

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"dblog/formatter"

	"dblog/trace"
)

// L is the global instance of the logger with the default configuration.
var L = New(zap.InfoLevel)

// Logger is a wrapper around the zap logger.
type Logger struct {
	*zap.SugaredLogger
}

// New returns a new instance of the logger with the specified log level.
func New(level zapcore.Level) *Logger {
	// This configuration shall be adjusted to the specific needs of the application.
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	return &Logger{zap.New(core).Sugar()}
}

// InfoSQLContext logs the operation, query, and arguments with the correlation ID extracted from the context.
// Provided query and arguments are formatted using the formatter package.
func (l *Logger) InfoSQLContext(ctx context.Context, operation, query string, args any) {
	l.SugaredLogger.Info(
		operation, " correlation_id=", trace.GetCorrelationID(ctx), formatter.QueryWithArgs(query, args),
	)
}

// InfoSQLContext logs the operation, query, and arguments with the correlation ID extracted from the context.
// Provided query and arguments are formatted using the formatter package.
func InfoSQLContext(ctx context.Context, operation, query string, args any) {
	L.InfoSQLContext(ctx, operation, query, args)
}

// InfoSQL logs the operation, query, and arguments.
// Provided query and arguments are formatted using the formatter package.
func (l *Logger) InfoSQL(operation, query string, args any) {
	l.SugaredLogger.Info(operation, formatter.QueryWithArgs(query, args))
}

// InfoSQL logs the operation, query, and arguments.
// Provided query and arguments are formatted using the formatter package.
func InfoSQL(operation, query string, args any) {
	L.InfoSQL(operation, query, args)
}
