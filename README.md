<div align="center" style="margin-bottom:20px">
  <!-- <img src=".assets/banner.png" alt="logger" /> -->
  <div align="center">
    <a href="https://github.com/blugnu/unilog/actions/workflows/qa.yml"><img alt="build-status" src="https://github.com/blugnu/unilog/actions/workflows/qa.yml/badge.svg?branch=master&style=flat-square"/></a>
    <a href="https://goreportcard.com/report/github.com/blugnu/unilog" ><img alt="go report" src="https://goreportcard.com/badge/github.com/blugnu/unilog"/></a>
    <a><img alt="go version >= 1.14" src="https://img.shields.io/github/go-mod/go-version/blugnu/unilog?style=flat-square"/></a>
    <a href="https://github.com/blugnu/unilog/blob/master/LICENSE"><img alt="MIT License" src="https://img.shields.io/github/license/blugnu/unilog?color=%234275f5&style=flat-square"/></a>
    <a href="https://coveralls.io/github/blugnu/unilog?branch=master"><img alt="coverage" src="https://img.shields.io/coveralls/github/blugnu/unilog?style=flat-square"/></a>
    <a href="https://pkg.go.dev/github.com/blugnu/unilog"><img alt="docs" src="https://pkg.go.dev/badge/github.com/blugnu/unilog"/></a>
  </div>
</div>

<br>

# unilog

> noun: def. a unified/universal logger

A package that provides an adaptable logger implementation to be used by other re-usable modules that wish to emit logs using a logger supplied by the consuming project, ensuring consistent log output from both the application and packages in any modules (that support `unilog`).

A module that supports logging via a `unilog.Logger` may also register enrichment functions to automatically add enrichment from context at the time of emitting any log entry.  Applications may also register their own enrichment functions and/or explicitly add enrichment to individual log entries.

## blugnu/errorcontext support

`unilog` supports [blugnu/errorcontext](https://github.com/blugnu/errorcontext) when logging messages using any of these functions:

* `Error()`
* `Errorf()`
* `FatalError()`
* `Fatalf()`
* `Warnf()`
* `Infof()`
* `Debugf()`
* `Tracef()`

For `Error()` and `FatalError()` the error being logged is checked for an `ErrorContext`.

For `Tracef()`, `Debugf()`, `Infof()`, `Warnf()`, `Errorf()` and `Fatalf()` args are inspected for any `error`s.  If an `error` is identified it is checked for an `ErrorContext`; if there is no `ErrorContext` further args are checked until an `ErrorContext` is found or there are no more args.

If an `ErrorContext` is identified, the context in the error is used to provide enrichment of the log entry before being emitted.

## How It Works

`unilog` does not implement an actual logger.  It provides a delegate that routes logging calls via an _adapter_ to a logger configured by an application.  A consuming project will configure whatever logger it wishes and then wrap that with the appropriate _adapter_ so that it may be injected into any modules or packages that support a `unilog.Logger`.

An adapter is provided for the standard library `log` package.  This may be initialised using `unitlog.StdLog()`.

A `Nul` adapter is also provided.  This produces no log output what-so-ever ("logging to NUL").

An adapter for [logrus](https://github.com/sirupsen/logrus) is available in a separate module: ([unilog4logrus](https://github.com/blugnu/unilog4logrus)).  The `logrus` adapter is provided in a separate module to avoid `unilog` itself taking any dependency on `logrus`.

<br>
<hr>
<br>

## How to Use UniLog

### In an Application

1. Configure a logger
2. Wrap your logger in a `unilog.Logger` using either `unilog.UsingAdapter()` or a helper func such as `unilog.Nul()` or `unilog.StdLog()`.  For any other adapters, refer to their documentation for any helper functions they may provide.
3. Pass your `unilog.Logger` into any modules/packages used that support it
4. _OPTIONAL:_ Register any `Enrichment` functions provided by your project
5. To emit logs, initialise an entry with any relevant context and emit messages as required
6. Enjoy your logs!

#### Example: Using unilog with std log

```golang
// Initialise a nul logger by default

import (
    "flag"

    "github.com/blugnu/unilog"

    "myorg/foo"
)

// By default logs will be "redirected to nul".  i.e. no log output 
var logger unilog.Logger = unilog.Nul()

func main() {
  // A hypothetical command flag to turn on logging.
  logEnabled := true
  flag.BoolVar(&logEnabled, "log", false, "turn on log output")
  flag.Parse()

  // If logging is enabled, replace the Nul() logger with a std log logger
  if logEnabled {
      logger = unilog.StdLog()
  }

  // Pass logger into the `foo` package, which supports unilog via an exported variable
  foo.Logger = logger

  // Do some logging ourselves
  log := logger.NewEntry()
  log.Info("logging initialised")

  // Any logs written by foo.SetupTheFoo() will use the same logger as 'log'
  if err := foo.SetupTheFoo(); err != nil {
    // If foo uses go-errorcontext to capture context with errors, the error
    // will automatically be logged with any enrichment available in the
    // captured context
    log.FatalError(err)
  }

  // ... etc
}
```

<br>

### In a Module/Package (to Support `unilog`)

1. Provide a mechanism for a `unilog.Logger` to be supplied for use by your module/package
2. _OPTIONAL_: Export any enrichment functions to be manually registered by projects
3. _OPTIONAL_: Implement an `init()` function to auto-register any enrichment functions
4. _OPTIONAL_: Add log enrichment data to `context` where appropriate and use `errorcontext` to return errors with context for enriched logging 
5. Write logs from your code using the supplied `Logger`
    - initialise a `unilog.Entry` in any function wishing to emit a log using either `FromContext()` or `NewEntry()` on the `unilog.Logger` (depending on whether a context is accepted by or otherwise available to the function)
    - emit logs at the appropriate level using the `Entry` obtained

> _**NOTE:** You should ensure that logs are **not written** if _no_ `Logger` is configured.</br></br>_**Either**_: ensure that logging statements are conditional (tedious and error prone)</br>_**or**_: initialise a default `unilog.Logger` using `unilog.Nul()`.</br></br>_**Alternatively** (recommended)_: treat the lack of a `Logger` as an error in any initialization provided by your module, requiring applications to _explicitly_ configure any `Logger`, including `Nul()`_.

### Implementing an Adapter

1. Implement the `unilog.Adapter` interface (see below)

2. _OPTIONAL_: provide an interface for any adapter specific configuration you wish to provide (the `Adapter` interface in [unilog4logrus](https://github.com/blugnu/unilog4logrus) is an example).   

3. _OPTIONAL (but recommended)_: Provide a helper function named `Logger` that accepts a logger and any additional required params that configures an adapter and returns the result of `unilog.UsingAdapter()`

### Adapter Interface

The following three functions are required to be implemented by an `Adapter`:

| function | description |
| -- | -- |
| `Emit(unilog.Level, string)` | implement this function to emit logs using whichever logging package your adapter supports.  Your adapter must map the `unilog.Level` to the corresponding level supported by the underlying logging package (or adapt the behaviour accordingly if the underlying logging package does not directly support leveled logging) |
| `NewEntry() Adapter` | implement this function to return a new adapter corresponding to a new log entry |
|	`WithField(string, any) Adapter` | implement this function to return a new adapter with the supplied, named value added to any log enrichment on the receiving adapter |

### Adapter Reference Example

The [unilog4logrus](https://github.com/unilog4logrus) adapter project provides a reference example, alongside the `Nul()` and `StdLog()` adapters implemented in the `unilog` package itself.