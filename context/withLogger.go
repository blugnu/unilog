package context

import (
	"context"

	"github.com/blugnu/unilog"
	"github.com/blugnu/unilog/internal"
)

type Context = context.Context

func WithLogger(ctx context.Context, log unilog.Logger) context.Context {
	return context.WithValue(ctx, internal.LoggerKey, log)
}
