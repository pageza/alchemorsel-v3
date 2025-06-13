package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingFunctions(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	SetLogger(zap.New(core).Sugar())

	Infof("hello %s", "world")
	Errorf("problem %d", 1)

	if recorded.Len() != 2 {
		t.Fatalf("expected 2 logs, got %d", recorded.Len())
	}
	if recorded.All()[0].Message != "hello world" {
		t.Fatalf("unexpected message: %s", recorded.All()[0].Message)
	}
	if recorded.All()[1].Level != zap.ErrorLevel {
		t.Fatalf("expected error level")
	}
}

func TestContextLogger(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	l := zap.New(core).Sugar()
	ctx := ToContext(context.Background(), l)
	FromContext(ctx).Info("msg")

	if recorded.Len() != 1 {
		t.Fatalf("expected 1 log, got %d", recorded.Len())
	}
}
