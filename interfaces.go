package unilog

import "context"

// Adapter is an interface the mediates between the Logger interface used
// by applications and modules and a concrete logger implementation.
type Adapter interface {
	Emit(Level, string)
	NewEntry() Adapter
	WithField(string, any) Adapter
}

// Logger is the interface used by applications and modules to initialise
// log entries.
type Logger interface {
	WithContext(context.Context) Entry // WithContext returns an Entry encapsulating the specific Context
	NewEntry() Entry                   // Returns a new Entry encapsulating the Context supplied when the Logger was initialised
}

// Enricher is an interface that provides a function that can add a named
// value (a 'Field') to a log entry.
//
// This interface separates this function from other functions supported
// by the Entry interface so that enrichment functions can receive a
// reference to the Entry as a more appropriate Enricher with only
// the enrichment function.
type Enricher interface {
	WithField(name string, value any) Entry // WithField returns a new adds a named value and returns the Entry
}

// Entry is the interface for an invidual log entry.  It provides functions
// for emitting messages at each of the logging levels supported by unilog.
//
// In addition, an Entry can be used to produce a new Entry encapsulating
// a new/different context, via the WithContext() function.
//
// Entry also includes the Enricher interface, so that an Entry may be
// enriched with one-off, specific fields if required, in addition to
// any enrichment provided automatically by any registered enrichment
// functions.
type Entry interface {
	Enricher
	Debug(s string)                    // Debug emits a Debug level log message
	Debugf(format string, args ...any) // Debugf emits a Debug level log message using a specified format string and args
	Error(err error)                   // Error emits an Error level log message consisting of err.Error()
	Errorf(format string, args ...any) // Errorf emits an Error level log message using a specified format string and args
	Fatal(s string)                    // Fatal emits a Fatal level log message then calls os.Exit(1)
	Fatalf(format string, args ...any) // Fatalf emits a Fatal level log message using a specified format string and args, then calls os.Exit(1)
	FatalError(err error)              // FatalError emits a Fatal level log message consisting of err.Error() then calls os.Exit(1)
	Info(s string)                     // Info emits an Info level log message
	Infof(format string, args ...any)  // Infof emits an Info level log message using a specified format string and args
	Trace(s string)                    // Trace emits a Trace level log message
	Tracef(format string, args ...any) // Tracef emits a Trace level log message using a specified format string and args
	Warn(s string)                     // Warn emits a Warn level log message
	Warnf(format string, args ...any)  // Warnf emits a Warn level log message using a specified format string and args
	WithContext(context.Context) Entry // WithContext returns a new Entry encapsulating the specified Context
}
