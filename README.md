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

A module that supports logging via a `unilog.Logger` may also register enrichment functions to automatically add enrichment from a context at the time of emitting a log entry.  Applications may also register their own enrichment functions.

## How It Works

`unilog` does not implement an actual logger.  It provides a delegate that routes logging calls via an _adapter_ to a logger configured by an application/service.  A consuming project will configure whatever logger it wishes and then wrap that with the appropriate _adapter_ so that it may be injected into any modules or packages that support a `unilog.Logger`.

Adapters are provided in separate modules for `logrus` ([unilog4logrus](https://github.com/unilog4logrus)) and the standard library `log` package ([unilog4stdlog](https://github.com/unilog4stdlog)).

A `Nul` adapter is provided by the `unilog` package itself; this adapter produces no log output what-so-ever ("logging to NUL").

<br>
<hr>
<br>

## How to Use UniLog

### Implementing an Adapter

1. Implement the `unilog.Adapter` interface
  - only 3 functions are required for each adapter
  - the [unilog4logrus](https://github.com/unilog4logrus) and [unilog4stdlog](https://github.com/unilog4stdlog) adapter projects provide reference examples
2. _OPTIONAL (but recommended)_: Provide a helper function named `Logger` that accepts a logger and any additional required params that configures an adapter and returns the result of `unilog.UsingAdapter()`

### In an Application or Service Project

1. Configure a logger
2. Wrap your logger in a `unilog.Logger` using either `unilog.UsingAdapter()` or a helper func provided by the adapter implementation (for example, `unilog4logrus.Logger()`)
4. Pass your `unilog.Logger` into any modules/packages used that support unilog
3. _OPTIONAL:_ Register any `Enrichment` functions provided by your project
5. Enjoy your logs!

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