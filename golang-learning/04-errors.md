# 04. Errors

Go has no exceptions. A function that can fail returns an `error` as its last
result, and the caller checks it. This is the single biggest adjustment for most
people, so this chapter is thorough.

## `error` is an interface

```go
type error interface {
	Error() string
}
```

Anything with an `Error() string` method is an error. `nil` means no error. That
is the whole mechanism.

`internal/quoterepo/model.go` defines one:

```go
type NoQuotesError struct {
}

func (e *NoQuotesError) Error() string {
	return "no quotes available"
}
```

and `repository.go` returns it:

```go
if quoteCount == 0 {
	err := &NoQuotesError{}
	log.Error("no quotes available", slog.Any("err", err))
	return "", err
}
```

## The check

```go
quote, err := quoteService.GetQuote(ctx)
if err != nil {
	return nil, err
}
```

You will write `if err != nil` hundreds of times. It is not noise; it is the
control flow of the program made visible. Every place an error can occur is
marked, and every place decides: handle it, return it, or deliberately ignore it
with `_`.

Rules of thumb:

- Check errors immediately after the call that produces them.
- Return early. The happy path stays at the left margin.
- Handle an error once. Either log it or return it, not both. (The repository
  logs in `Gimme` and also returns; that is a judgement call for a leaf function
  that knows the most context. Do not do it at every layer.)
- Never ignore an error silently unless you can say why. `_ = json.Unmarshal(data, quotes)`
  in `NewQuoteRepository` is a deliberate choice: embedded data is validated at
  build time, so an empty repository is an acceptable outcome that `Gimme` reports
  later. Write a comment when you do this.

## Creating errors

```go
errors.New("Server context must be initialized before start")   // fixed message
fmt.Errorf("reading %s: %w", path, err)                           // formatted, wrapping err
```

Use `errors.New` for constants, `fmt.Errorf` when you need to include values.
`%w` (not `%v`) wraps the inner error so that it can be inspected later.

Define a type when callers need to distinguish the error or read data from it,
as `NoQuotesError` does. Define a sentinel variable when they only need
identity:

```go
var ErrNotFound = errors.New("not found")
```

## Wrapping and inspecting

Wrapping adds context as an error travels up the stack:

```go
if err := repo.Load(); err != nil {
	return fmt.Errorf("loading quotes: %w", err)
}
```

The message reads like a path: `loading quotes: open quotes.json: no such file`.
Do not start messages with capitals or end with punctuation, because they get
joined with `: `.

To ask "is this, at any depth, a particular error?":

```go
if errors.Is(err, ErrNotFound) { ... }          // sentinel identity

var noQuotes *quoterepo.NoQuotesError
if errors.As(err, &noQuotes) { ... }            // typed, and you get the value
```

`errors.Is` and `errors.As` walk the wrap chain, so they work no matter how many
layers added context. Comparing with `==` or a type switch only sees the outer
layer, so prefer the `errors` functions.

## Errors at the boundary: Twirp

Errors that cross the network need a code the client understands. Twirp defines
`twirp.Error`, an interface with `Code()`, `Msg()` and `Meta()` in addition to
`Error()`. The service layer is where plain errors become Twirp errors:

```go
func (cs *QuoteService) GetQuote(ctx context.Context) (string, twirp.Error) {
	quote, err := cs.repo.Gimme(ctx)
	if err != nil {
		return "", twirp.InternalError(err.Error())
	}
	return quote, nil
}
```

and validation produces a client error:

```go
return twirp.InvalidArgumentError("userSession.userUrn", "invalid user urn")
```

The generated server turns these into JSON with the right HTTP status. A plain
`error` returned from a handler becomes a generic 500 `internal` error, so
translating at the service layer is what gives clients precise codes.

Note the signature returns `twirp.Error`, a concrete interface, rather than
`error`. That is fine here because every caller is a Twirp handler. It does
carry one trap, covered next.

## The nil interface trap

An interface value is `nil` only if both its type and value are `nil`. This
compiles and is wrong:

```go
func find() error {
	var e *NoQuotesError = nil
	return e            // returns a non-nil error holding a nil pointer
}
if err := find(); err != nil { /* this runs */ }
```

Return a literal `nil` from functions whose result type is an interface, never a
typed nil pointer. `validateUserSession` in `internal/api/validation.go` does
this correctly:

```go
func validateUserSession(session *common.UserSession) twirp.Error {
	if ... {
		return twirp.InvalidArgumentError(...)
	}
	return nil
}
```

## Panics

`panic` unwinds the stack and, if nothing `recover`s, crashes the program with a
trace. It is for bugs and unrecoverable start-up failures, not for expected
conditions. The repository panics exactly where that is appropriate:

```go
if sctx == nil {
	err := errors.New("Server context must be initialized before start")
	log.Error("Failed to start server", slog.Any("err", err))
	panic(err)
}
```

and

```go
if err := server.ListenAndServe(); err != http.ErrServerClosed {
	panic(err)
}
```

A server that cannot bind its port has nothing better to do than exit loudly.
`log.Fatal` and `log.Fatalf` print and call `os.Exit(1)` without unwinding, and
`context.go` uses them for the same class of failure.

Do not `recover` to turn panics into errors as a general pattern. The generated
Twirp server already recovers from panics inside handlers and answers with an
`internal` error, so a bug in one request does not take down the process.

## Try it

1. Change `Gimme` to return `fmt.Errorf("picking a quote: %w", err)` and update
   the service to use `errors.As` to detect `*NoQuotesError` and return
   `twirp.NotFoundError` for that case only. Run the unit tests.
2. Make `NoQuotesError` carry the number of quotes it expected and include it in
   the message.
3. Write a function that deliberately returns a typed nil pointer as `error`,
   check it with `!= nil`, and print `fmt.Printf("%T %v\n", err, err)` to see the
   type and value.
4. Read `go doc errors` and `go doc twirp.Error` in the terminal.

## Read next

- [Go blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- [Twirp errors](https://twitchtv.github.io/twirp/docs/errors.html)
- [05. Structs, methods and interfaces](05-structs-methods-interfaces.md)
