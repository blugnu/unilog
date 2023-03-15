package unilog

import "context"

type Adapter interface {
	Emit(Level, string)
	NewEntry() Adapter
	WithField(string, any) Adapter
}

type Logger interface {
	WithContext(context.Context) Entry
	NewEntry() Entry
}

type Enricher interface {
	WithField(name string, value any) Entry
}

type Entry interface {
	Enricher
	Debug(s string)
	Debugf(format string, args ...any)
	Error(err error)
	Errorf(format string, args ...any)
	Fatal(s string)
	Fatalf(format string, args ...any)
	FatalError(err error)
	Info(s string)
	Infof(format string, args ...any)
	Trace(s string)
	Tracef(format string, args ...any)
	Warn(s string)
	Warnf(format string, args ...any)
	WithContext(context.Context) Entry
}
