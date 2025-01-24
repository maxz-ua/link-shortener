package slogdiscard_test

import (
	"context"
	"link-shortener/internal/lib/logger/handlers/slogdiscard"
	"log/slog"
	"testing"
)

func TestDiscardHandler(t *testing.T) {
	handler := slogdiscard.NewDiscardHandler()
	logger := slog.New(handler)

	// Test Enabled
	if handler.Enabled(context.Background(), slog.LevelInfo) {
		t.Errorf("Expected Enabled to return false")
	}
	// Test Handle
	err := handler.Handle(context.Background(), slog.Record{})
	if err != nil {
		t.Errorf("Handle should return nil, got: %v", err)
	}
	// Test WithAttrs and WithGroup
	if handler.WithAttrs(nil) != handler {
		t.Errorf("WithAttrs should return the same instance")
	}
	if handler.WithGroup("test") != handler {
		t.Errorf("WithGroup should return the same instance")
	}
	// Test logger does nothing
	logger.Info("This should not appear anywhere")
}
