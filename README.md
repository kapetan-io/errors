# Errors
Is a drop in replacement library for std "errors" package with support for attaching context logging fields for 
use with `log/slog`. This is to facilitate the idiom [only handle errors once rule](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
while not losing any context where the error occurred.

## Examples
Attach context to error messages
```go
// Attach context to an error
err := errors.With("foo", "bar").Error("query failed")

// Prints `query failed`
fmt.Printf("%s\n", err)

// Prints `query failed (foo=bar)`
fmt.Printf("%+v\n", err)

// Prints `[foo=bar]`
fmt.Printf("%v\n", errors.AttrsFrom(err))

// Prints `2024/09/30 12:00:00 level=ERROR msg="query failed" foo=bar`
slog.LogAttrs(ctx, slog.LevelError, err.Error(), errors.AttrsFrom(err)...)
```
Use with `slog.Attr` to create complex logging context for packages or structs
```go
s := Struct{
   attrs: errors.WithAttr(slog.String("database", "sqlite")),
}

// Prints `query failed (database=sqlite, foo=bar)`
fmt.Printf("%+v\n", s.attrs.With("foo", "bar").Error("query failed"))
```
Extract OTEL standard `code` location information for use with slog
```go
// Prints `Attributes [
//   code.function=github.com/kapetan-io/errors_test.ExampleAttrs
//   code.filepath=.../errors/example_test.go
//   code.lineno=16
//   foo=bar]
fmt.Printf("Attributes %v\n", errors.AttrsWithCodeLoc(err))

// 2024/09/30 12:31:16 ERROR query failed foo=bar 
// code.filepath=/Users/thrawn/Development/errors/example_test.go
// code.function=github.com/kapetan-io/errors_test.ExampleAttrs code.lineno=16
slog.LogAttrs(ctx, slog.LevelError, err.Error(), errors.AttrsWithCodeLoc(err)...)
```
Works with standard golang error wrapping
```go
err = errors.New("query error")
wrap := errors.With("key1", "value1").Errorf("message: %w", err)

// Prints `message: query error (key1=value1)`
fmt.Printf("%+v\n", wrap)
```
Use standard introspection functions to extract fields
```go
var f errors.HasAttrs
if errors.As(wrap, &f) {
    slog.LogAttrs(ctx, slog.LevelError, "this error has attrs in the chain",
       errors.AttrsFrom(err)...)
}
```

## Include pass through std 'error' library methods
Provides pass through access to the standard `errors.Is()`, `errors.As()`, `errors.Unwrap()` so you don't need to
import this package and the standard `errors` package.


### API
- **errors.With()** - Attach context to an error in the form of key value pairs `errors.With("key", "value")`
- **errors.WithAttr()** - Attach context to an error using `slog.Attr`  `errros.WithAttr(slog.String("key", "value"))`
- **errors.With().Wrap()** - Wrap an error without a message, attaching the code location where `Wrap()` was called
- **errors.Last()** - Same as standard lib `errors.As()` but returns the last error in the err tree instead
- **errors.With().Error()** - Same as standard lib `errors.New()` includes code location where `Error()` was called
- **errors.With().Errorf()** - Same as standard lib `fmt.Errorf()` includes code location where `Errorf()` was called
- **errors.Wrap()** - Wrap an error without a message, including the code location where `Wrap()` was called
- **errors.Unwrap()** - Same as standard lib `errors.Unwrap()`
- **errors.Error()** - Same as standard lib `errors.New()` includes code location where `Error()` was called
- **errors.Errorf()** - Same as standard lib `fmt.Errorf()` include code location where `Errorf()` was called
- **errors.Wrap()** - Wrap an error without a message, including the code location where `Wrap()` was called
- **errors.New()** - Same as standard lib `errors.New()`
- **errors.As()** - Same as standard lib `errors.As()`
- **errors.Is()** - Same as standard lib `errors.Is()`
  of the first.
