package unilog

import (
	"context"
	"fmt"

	"github.com/blugnu/go-errorcontext"
)

// logger implements `Logger`, `Enricher` and `Entry` interfaces.  It encapsulates
// a specific context with an `Adapter` which is used to perform all concrete
// logging operations.
type logger struct {
	context.Context
	Adapter
}

// emit sends a specified string to the logger with the specified log level.
func (log *logger) emit(level Level, s string) {
	entry := log.fromContext(log.Context)
	entry.Emit(level, s)
}

// fromContext returns a new `logger` using the same `Adapter` as the receiver,
// encapsulating the specified `Context`.  The new `logger` has all registered
// enrichment applied.
func (log *logger) fromContext(ctx context.Context) *logger {
	logger := &logger{ctx, log.Adapter.NewEntry()}

	var enriched Enricher = logger
	for _, enrich := range enrichmentFuncs {
		enriched = enrich(ctx, enriched)
	}

	return logger
}

// Trace emits a string as a `Trace` level entry to the log.
func (log *logger) Trace(s string) {
	log.emit(Trace, s)
}

// Tracef emits a `Trace` level entry to the log using a format string and args.
func (log *logger) Tracef(format string, args ...any) {
	log.Trace(fmt.Sprintf(format, args...))
}

// Debug emits a string as a `Debug` level entry to the log.
func (log *logger) Debug(s string) {
	log.emit(Debug, s)
}

// Debugf emits a `Debug` level entry to the log using a format string and args.
func (log *logger) Debugf(format string, args ...any) {
	log.Debug(fmt.Sprintf(format, args...))
}

// Info emits a string as an `Info` level entry to the log.
func (log *logger) Info(s string) {
	log.emit(Info, s)
}

// Infof emits an `Info` level entry to the log using a format string and args.
func (log *logger) Infof(format string, args ...any) {
	log.Info(fmt.Sprintf(format, args...))
}

// Warn emits a string as a `Warn` level entry to the log.
func (log *logger) Warn(s string) {
	log.emit(Warn, s)
}

// Warnf emits a `Warn` level entry to the log using a format string and args.
func (log *logger) Warnf(format string, args ...any) {
	log.Warn(fmt.Sprintf(format, args...))
}

// Error emits an error as an `Error` level entry to the log.
//
// If the error wraps a specific context then the error is logged using an entry
// enriched with any  information in the context supported by a registered
// enrichment function.
func (log *logger) Error(err error) {
	ctx := errorcontext.FromError(log.Context, err)
	entry := log.fromContext(ctx)
	entry.emit(Error, err.Error())
}

// Errorf emits an `Error` level entry to the log using a format string and args.
//
// If the error wraps a specific context then the error is logged using an entry
// enriched with any  information in the context supported by a registered
// enrichment function.
func (log *logger) Errorf(format string, args ...any) {
	entry := log
	for _, a := range args {
		if err, isError := a.(error); isError {
			ctx := errorcontext.FromError(log.Context, err)
			entry = entry.fromContext(ctx)
			break
		}
	}
	entry.Error(fmt.Errorf(format, args...))
}

// Fatal emits a string as a `Fatal` level entry to the log then terminates
// the process with an exit code of 1.
func (log *logger) Fatal(s string) {
	log.emit(Fatal, s)
	exit(1)
}

// Fatalf emits a `Fatal` level entry to the log using a format string and args
// then terminates the process with an exit code of 1..
func (log *logger) Fatalf(format string, args ...any) {
	entry := log
	for _, a := range args {
		if err, isError := a.(error); isError {
			ctx := errorcontext.FromError(log.Context, err)
			entry = entry.fromContext(ctx)
			break
		}
	}
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
	entry := log.Adapter.WithField(name, value)
	return &logger{log.Context, entry}
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
	return &logger{ctx, adapter}
}
