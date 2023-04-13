package unilog

import (
	"context"
	"testing"
)

func TestContextWithLogger(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	logger := &logger{}

	// ACT
	ctx = ContextWithLogger(ctx, logger)

	// ASSERT
	got := ctx.Value(loggerContextKey)
	wanted := logger
	if wanted != got {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}

func TestLogFromContext(t *testing.T) {
	t.Run("when context does not contain a Logger", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()

		// ACT
		entry := LogFromContext(ctx)

		// ASSERT
		wanted := (Entry)(nil)
		got := entry
		if wanted != got {
			t.Errorf("wanted %#v, got %#v", wanted, got)
		}
	})

	t.Run("when context contains a Logger", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()
		ctx = context.WithValue(ctx, loggerContextKey, nul)

		// ACT
		result := LogFromContext(ctx)

		// ASSERT
		wanted := nul.WithContext(ctx).(*logger)
		got := result.(*logger)
		if *wanted != *got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}

func TestLoggerFromContext(t *testing.T) {
	t.Run("when context does not contain a Logger", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()

		// ACT
		result := LoggerFromContext(ctx)

		// ASSERT
		wanted := (Logger)(nil)
		got := result
		if wanted != got {
			t.Errorf("wanted %#v, got %#v", wanted, got)
		}
	})

	t.Run("when context contains a Logger", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()
		ctx = context.WithValue(ctx, loggerContextKey, nul)

		// ACT
		result := LoggerFromContext(ctx)

		// ASSERT
		wanted := nul
		got := result.(*logger)
		if *wanted != *got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}
