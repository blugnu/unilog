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

> noun: def. a uni[_fied_] log[_ger_]

A package that provides an adaptable logger implementation to be used by other re-usable modules that wish to emit logs using a logger supplied by the consuming project, ensuring consistent log output from both the application and pacakges in any modules (that support `unilog`).

A module that supports logging via a `unilog.Logger` may also register enrichment functions to automatically add enrichment from context at the time of emitting any log entry.  Applications may also register their own enrichment functions and/or explicitly add enrichment to individual log entries.

## How It Works

`unilog` does not implement an actual logger.  It provides a delegate that routes logging calls via an _adapter_ to a logger configured by an application/service.  A consuming project will configure whatever logger it wishes and then wrap that with the appropriate _adapter_ so that it may be injected into any modules or packages that support a `unilog.Logger`.

An adapter is provided for the standard library `log` package. A `Nul` adapter is also provided.  This produces no log output what-so-ever ("logging to NUL").

An adapter for logrus is available in a separate module: ([unilog4logrus](https://github.com/unilog4logrus)).  Providing this adapter in a separate module avoids unilog itself taking any dependency on `logrus` itself, for those using alternate logging packages.

<br>
<hr>
<br>

## How to Use UniLog

### In an Application or Service Project

1. Configure a logger
2. Wrap your logger in a `unilog.Logger` using either `unilog.UsingAdapter()` or a helper func provided by the adapter implementation (for example, `unilog4logrus.Logger()`)
4. Pass your `unilog.Logger` into any modules/packages used that support unilog
3. _OPTIONAL:_ Register any `Enrichment` functions provided by your project
5. Enjoy your logs!

#### Example: Using unilog with logrus

```golang
var logger unilog.Logger

func main() {
  // A hypothetical command flag to suppress all log output, otherwise
  // we will log using std log
  logEnabled := true
  flag.BoolVar(&noLog, "nolog", false, "suppress all log output")
	flag.Parse()

  // Get a unilog encapsulating std log
	logger = unilog.StdLog()

  // Pass logger into the `foo` module, which supports unilog
  foo.Logger = logger

  // Do some logging ourselves
  log := logger.NewEntry()
  log.Info("logging initialised")

  // Any logs written by SetupTheFoo() will use the same logger as 'log'
  if err := foo.SetupTheFoo(); err != nil {
    log.FatalError(err)
  }

  // ... etc
}
```

<br>

### In a Module/Package (to Support `unilog`)

1. Provide a mechanism for a `unilog.Logger` to be supplied for use by your module/package
2. _OPTIONAL_: Export an enrichment function (to be manually registered by projects if they wish)
3. _OPTIONAL_: Implement an `init()` function to auto-register any enrichment functions
4. _OPTIONAL_: Add log enrichment data to `context` where appropriate and use `errorcontext` to return errors with context for enriched logging 
4. Write logs from your code using the supplied `Logger`
    - initialise a `unilog.Entry` in any function wishing to emit a log using either `FromContext()` or `NewEntry()` on the `unilog.Logger` (depending on whether a context is accepted by or otherwise available to the function)
    - emit logs at the appropriate level using the `Entry` obtained
    - Ensure that logs are **not written** if _no_ `Logger` is configured.  _Either_ ensure logging statements are condition _or_ initialise a default `unilog.Logger` using `unilog.Nul()`

### Implementing an Adapter

1. Implement the `unilog.Adapter` interface
  - only 3 functions are required for an adapter
  - the [unilog4logrus](https://github.com/unilog4logrus) adapter project provides a reference example, alongside the Nul and StdLog adapters implemented in the unilog package itself.
2. _OPTIONAL_: provide an interface for any adapter specific configuration you wish to provide (the `Adapter` interface in [unilog4logrus](https://github.com/unilog4logrus) is an example).   
3. _OPTIONAL (but recommended)_: Provide a helper function named `Logger` that accepts a logger and any additional required params that configures an adapter and returns the result of `unilog.UsingAdapter()`

