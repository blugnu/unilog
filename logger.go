package unilog

import (
	"context"
	"fmt"

	"github.com/blugnu/go-errorcontext"
)

// logger implements Logger, Enricher and Entry interfaces.  It encapsulates
// a specific context with an Adapter which is used to perform all concrete
// logging operations.
type logger struct {
	context.Context
	Adapter
	fields map[string]any
}

// copyFields returns a copy of the fields map, or nil if there are no fields.
func (log *logger) copyFields() map[string]any {
	if log.fields == nil {
		return nil
	}

	fields := make(map[string]any, len(log.fields))
	for k, v := range log.fields {
		fields[k] = v
	}
	return fields
}

// Emit sends a specified string to the logger with the specified log level.
func (log *logger) Emit(level Level, s string) {
	adapter := log.fromContext(log.Context).(*logger).Adapter

	for k, v := range log.fields {
		adapter = adapter.WithField(k, v)
	}

	adapter.Emit(level, s)
}

// entryFromArgs examines args to identify any error values.  If any error
// is found that contains a context (wrapped in an ErrorContext) then an
// entry is initialised with the context in the first such error, and returned.
//
// Otherwise, if there is no `error` in the args, or no `error`
// wrapping a context with an `ErrorContext`, the function simply returns
// a reference to the current entry.
func (log *logger) entryFromArgs(args ...any) Entry {
	for _, a := range args {
		if _, isError := a.(error); !isError {
			continue
		}

		ctx := errorcontext.FromError(log.Context, a.(error))
		if ctx == log.Context {
			continue
		}

		return log.fromContext(ctx)
	}

	return log
}

// fromContext returns a new `logger` using the same `Adapter` as the receiver,
// encapsulating the specified `Context`.  The new `logger` has all registered
// enrichment applied.
func (log *logger) fromContext(ctx context.Context) Entry {
	var logger = &logger{ctx, log.Adapter, log.copyFields()}

	var enriched Entry = logger
	for _, enrich := range enrichmentFuncs {
		enriched = enrich(ctx, enriched)
	}

	return enriched
}

// Trace emits a string as a `Trace` level entry to the log.
func (log *logger) Trace(s string) {
	log.Emit(Trace, s)
}

// Tracef emits a `Trace` level entry to the log using a format string and args.
func (log *logger) Tracef(format string, args ...any) {
	entry := log.entryFromArgs(args...)
	entry.Trace(fmt.Sprintf(format, args...))
}

// Debug emits a string as a `Debug` level entry to the log.
func (log *logger) Debug(s string) {
	log.Emit(Debug, s)
}

// Debugf emits a `Debug` level entry to the log using a format string and args.
func (log *logger) Debugf(format string, args ...any) {
	entry := log.entryFromArgs(args...)
	entry.Debug(fmt.Sprintf(format, args...))
}

// Info emits a string as an `Info` level entry to the log.
func (log *logger) Info(s string) {
	log.Emit(Info, s)
}

// Infof emits an `Info` level entry to the log using a format string and args.
func (log *logger) Infof(format string, args ...any) {
	entry := log.entryFromArgs(args...)
	entry.Info(fmt.Sprintf(format, args...))
}

// Warn emits a string as a `Warn` level entry to the log.
func (log *logger) Warn(s string) {
	log.Emit(Warn, s)
}

// Warnf emits a `Warn` level entry to the log using a format string and args.
func (log *logger) Warnf(format string, args ...any) {
	entry := log.entryFromArgs(args...)
	entry.Warn(fmt.Sprintf(format, args...))
}

// Error emits an error or string as an Error level entry to the log.
//
// If logging an error wrapping a specific context then the error is
// logged using an entry  enriched with any information in the error context.
//
// If logging a string then the string is logged as-is.
//
// If logging any other type then the type is logged using `fmt.Sprintf` with
// the `%v` format.
func (log *logger) Error(err any) {
	switch err := err.(type) {
	case error:
		ctx := errorcontext.FromError(log.Context, err)
		entry := log.fromContext(ctx)
		entry.Emit(Error, err.Error())
	case string:
		log.Emit(Error, err)
	default:
		log.Emit(Error, fmt.Sprintf("%v", err))
	}
}

// Errorf emits an `Error` level entry to the log using a format string and args.
//
// If the error wraps a specific context then the error is logged using an entry
// enriched with any  information in the context supported by a registered
// enrichment function.
func (log *logger) Errorf(format string, args ...any) {
	entry := log.entryFromArgs(args...)
	entry.Error(fmt.Errorf(format, args...))
}

// Fatal emits a string as a `Fatal` level entry to the log then terminates
// the process with an exit code of 1.
func (log *logger) Fatal(s string) {
	log.Emit(Fatal, s)
	exit(1)
}

// Fatalf emits a `Fatal` level entry to the log using a format string and args
// then terminates the process with an exit code of 1..
func (log *logger) Fatalf(format string, args ...any) {
	entry := log.entryFromArgs(args...)
	entry.Fatal(fmt.Sprintf(format, args...))
}

// FatalError emits an error as a `Fatal` level entry to the log.
//
// If the error wraps a specific context then the error is logged using an entry
// enriched with any  information in the context supported by a registered
// enrichment function.
func (log *logger) FatalError(err error) {
	ctx := errorcontext.FromError(log.Context, err)
	entry := log.fromContext(ctx)
	entry.Fatal(err.Error())
}

// WithField returns a new `Entry` enriched with an additional
// named field with the specified value.
func (log *logger) WithField(name string, value any) Entry {
	fields := log.copyFields()
	if fields == nil {
		fields = map[string]any{}
	}
	fields[name] = value

	return &logger{log.Context, log.Adapter, fields}
}

// WithContext returns a new `Entry`, enriched with any information
// available from a supplied context `ctx`.
func (log *logger) WithContext(ctx context.Context) Entry {
	return log.fromContext(ctx)
}

// NewEntry returns a new `Entry` enriched with any information
// available from the context established in the current entry
// (i.e. the receiver)
//
// Given:
//
//	e1 := log.FromContext(ctx)
//	e2 := log.FromContext(ctx).NewEntry()
//
// `e1“ and `e2“ are distinct and separate `LogEntry` values,
// both enriched with information available from the context `ctx`.
func (log *logger) NewEntry() Entry {
	return log.fromContext(log.Context)
}

// UsingAdapter initialises a new Logger encapsulating a specified
// context and using a supplied `Adapter`.
func UsingAdapter(ctx context.Context, adapter Adapter) Logger {
	return &logger{ctx, adapter, map[string]any{}}
}
