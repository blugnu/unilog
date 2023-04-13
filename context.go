package unilog

import (
	"context"
)

type contextKey int

const loggerContextKey contextKey = iota

// ContextWithLogger adds a Logger reference to a parent context.  The new context
// containing the Logger is returned.
//
// This function is intended to be used in conjunction with modules that support
// unilog logging by accepting a Logger passed in within a Context, rather than any
// more explicit and direct configuration value.
func ContextWithLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, log)
}

// LogFromContext inspects a specified context for a `unilog.Logger`;
// if a `Logger` is found, it is used to initialise a new `Entry` from
// the context which is then returned.
//
// If the context does not contain a `Logger`, `nil` is returned.
//
// NOTE: This function is intended to be used in modules that choose to accept
// a `Logger` supplied via a context, rather than providing a specific
// configuration variable or field.  This may be desirable where `Logger`
// support is added without wishing to break existing configuration contracts.
func LogFromContext(ctx context.Context) Entry {
	log := ctx.Value(loggerContextKey)
	if log == nil {
		return nil
	}

	return log.(Logger).WithContext(ctx)
}

// LoggerFromContext inspects a specified context for a `unilog.Logger`;
// if a `Logger` is found it is returned.
//
// If the context does not contain a `Logger`, `nil` is returned.
//
// NOTE: This function is intended to be used in modules that choose to accept
// a `Logger` supplied via a context, rather than providing a specific
// configuration variable or field.  This may be desirable where `Logger`
// support is added without wishing to break existing configuration contracts.
func LoggerFromContext(ctx context.Context) Logger {
	log := ctx.Value(loggerContextKey)
	if log == nil {
		return nil
	}

	return log.(Logger)
}
