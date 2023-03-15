package unilog

import "os"

// ExitFn is a func var allowing `unilog` code paths that reach an `os.Exit` call
// to be replaced by a non-exiting behaviour.
//
// This is primarily used in `unilog` unit tests but also allows (e.g.) a
// microservice to intercept such exit calls and perform a controlled exit,
// even when a log call results in termination of the microservice.
var ExitFn func(int) = os.Exit

// exit calls the `ExitFn` with the specified exit code.  Code paths in `unilog`
// that require termination of the process (e.g. `log.FatalError()`) call this
// `exit` function which in turn calls the `ExitFn` func var.
//
// To prevent `unilog` causing a process to terminate, replace `ExitFn`.
func exit(code int) {
	ExitFn(code)
}
