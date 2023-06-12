package unilog

import "context"

// Emitter is an interface that provides a single function for emitting
// a log message at a specified level.
type Emitter interface {
	Emit(Level, string)
}

// Enricher is in interface that provides a single function for enriching
// a log entry with an additional field having a specified name and value.
type Enricher interface {
	WithField(string, any) Entry
}

// Adapter is an interface that mediates between the abstract Logger and
// a concrete implementation using some logging package.
//
// An Adapter is also an Emitter.
type Adapter interface {
	Emitter
	NewEntry() Adapter
	WithField(string, any) Adapter
}

// Logger is the interface used by applications and modules to initialise
// log entries.  Applications should normally initialise a Logger with a
// desired Adapter, passing the Logger to packages that support unilog.
type Logger interface {
	WithContext(context.Context) Entry // WithContext returns an Entry encapsulating the specific Context
	NewEntry() Entry                   // Returns a new Entry encapsulating the Context supplied when the Logger was initialised
}

// Entry is the interface for an individual log entry.  An Entry is an Emitter
// that additionally provides helper functions for emitting messages at each of
// the supported logging levels.
//
// In addition, an Entry can be used to produce a new Entry encapsulating
// a new context, via the WithContext() function.
//
// Entry also provides a WithField function for providing one-off enrichment
// of individual (or related) entries, in addition to any enrichment provided
// from the logging context by registered enrichment functions.
type Entry interface {
	Emitter
	Debug(s string)                         // Debug emits a Debug level log message
	Debugf(format string, args ...any)      // Debugf emits a Debug level log message using a specified format string and args
	Error(err any)                          // Error emits an Error level log message consisting of err
	Errorf(format string, args ...any)      // Errorf emits an Error level log message using a specified format string and args
	Fatal(s string)                         // Fatal emits a Fatal level log message then calls os.Exit(1)
	Fatalf(format string, args ...any)      // Fatalf emits a Fatal level log message using a specified format string and args, then calls os.Exit(1)
	FatalError(err error)                   // FatalError emits a Fatal level log message consisting of err.Error() then calls os.Exit(1)
	Info(s string)                          // Info emits an Info level log message
	Infof(format string, args ...any)       // Infof emits an Info level log message using a specified format string and args
	Trace(s string)                         // Trace emits a Trace level log message
	Tracef(format string, args ...any)      // Tracef emits a Trace level log message using a specified format string and args
	Warn(s string)                          // Warn emits a Warn level log message
	Warnf(format string, args ...any)       // Warnf emits a Warn level log message using a specified format string and args
	WithContext(context.Context) Entry      // WithContext returns a new Entry encapsulating the specified Context
	WithField(name string, value any) Entry // WithField returns a new Entry with the named value added (a one-off enrichment)
}
