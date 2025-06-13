package logger

import (
	"context"
	"go.uber.org/zap"
)

// ctxKey is the type used for storing the logger in a context.Context.
type ctxKey struct{}

var (
	key  ctxKey
	base *zap.SugaredLogger
)

func init() {
	Init()
}

// Init initializes the base logger.
func Init() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	base = l.Sugar()
}

// SetLogger replaces the underlying logger. Mainly used for tests.
func SetLogger(l *zap.SugaredLogger) {
	base = l
}

// Logger returns the base logger instance.
func Logger() *zap.SugaredLogger {
	return base
}

// FromContext returns the request-scoped logger stored in the context. If none
// is found, the base logger is returned.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return base
	}
	if l, ok := ctx.Value(key).(*zap.SugaredLogger); ok {
		return l
	}
	return base
}

// ToContext stores the given logger in the context and returns the new context.
func ToContext(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, key, l)
}

// Infof logs an informational message using the base logger.
func Infof(format string, args ...any) {
	base.Infof(format, args...)
}

// Errorf logs an error message using the base logger.
func Errorf(format string, args ...any) {
	base.Errorf(format, args...)
}

// Fatal logs a fatal error message and exits.
func Fatal(err error) {
	base.Fatal(err)
}
