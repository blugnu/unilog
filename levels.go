package unilog

// Level identifies the logging level for a particular log entry.
// The possible values for `Level` are modelled on `logrus`,
// though `Panic` is not supported (`Fatal` is the most severe).
type Level int

const (
	Fatal Level = iota // logging at this level will terminate the process after emitting the log entry (without returning any error)
	Error
	Warn
	Info
	Debug
	Trace
)
