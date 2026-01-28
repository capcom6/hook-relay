package watermillfx_test

import (
	"testing"

	"github.com/capcom6/hook-relay/pkg/watermillfx"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	// Create a zap logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create zap logger: %v", err)
	}
	defer logger.Sync()

	// Create our watermill logger
	watermillLogger := watermillfx.NewLogger(logger)

	// Test each method
	testFields := watermill.LogFields{
		"test":   "value",
		"number": 42,
	}

	// Test Debug
	watermillLogger.Debug("test debug message", testFields)

	// Test Info
	watermillLogger.Info("test info message", testFields)

	// Test Trace
	watermillLogger.Trace("test trace message", testFields)

	// Test Error
	testErr := assert.AnError
	watermillLogger.Error("test error message", testErr, testFields)

	// Test With
	withLogger := watermillLogger.With(watermill.LogFields{"additional": "field"})
	withLogger.Info("test with logger", testFields)

	// Verify that the logger implements the LoggerAdapter interface
	var _ watermill.LoggerAdapter = watermillLogger
}
