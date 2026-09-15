# Coming from another language

A cheat sheet of Go equivalents for things you already know, and the habits to
unlearn. Skim the column for your background, then the "unlearn" list at the
end, which applies to everyone.

## Equivalents

| You know                          | In Go                                                                 |
| --------------------------------- | --------------------------------------------------------------------- |
| Class with fields and methods     | `struct` plus methods with a receiver ([05](05-structs-methods-interfaces.md)) |
| Constructor                       | A function `NewT(...) *T` by convention                               |
| Inheritance                       | Embedding for reuse; interfaces for polymorphism. No hierarchy         |
| `interface` / `implements`        | Interface satisfied implicitly by having the methods                  |
| Abstract class                    | Interface plus a struct with shared behaviour you embed                |
| `public` / `private`              | Capitalised = exported from the package; lower case = package only     |
| `this` / `self`                   | The receiver, named by you (`cs`, `h`, `cr`)                          |
| Exceptions, `try` / `catch`       | Return `error`; check `if err != nil` ([04](04-errors.md))            |
| `finally` / `using` / `with`      | `defer`                                                               |
| `null` / `None` / `nil`           | `nil`, valid for pointers, slices, maps, channels, interfaces, funcs  |
| Optional / nullable scalar        | Pointer to the type, or a wrapper message in protobuf                 |
| Generics `<T>`                    | Type parameters `[T any]`, since 1.18; use sparingly                   |
| Enum                              | Typed constants with `iota`                                           |
| Ternary `a ? b : c`               | `if` / `else`; there is no ternary                                    |
| `List<T>` / array                 | Slice `[]T`; `append` returns the new slice                            |
| `Map<K,V>` / dict                 | `map[K]V`; comma-ok read `v, ok := m[k]`                               |
| `for each`                        | `for i, v := range xs`                                                |
| `while`                           | `for cond { }`                                                        |
| Lambda / closure                  | `func(x int) int { return x * 2 }`, captures by reference             |
| Thread / async task               | `go f()` goroutine ([06](06-context-goroutines-shutdown.md))           |
| `async` / `await`, promises       | Goroutines plus channels; blocking calls are fine and cheap            |
| Queue between threads             | `chan T`                                                              |
| Mutex                             | `sync.Mutex`, `sync.RWMutex`                                          |
| Cancellation token                | `context.Context`                                                     |
| Decorator / annotation for data   | Struct tags `` `json:"name"` ``                                       |
| Middleware                        | `func(http.Handler) http.Handler`                                     |
| Package manager + lockfile        | `go.mod` + `go.sum`, built in ([01](01-toolchain-and-modules.md))      |
| Formatter / linter config         | `gofmt` has no options; `golangci-lint` for the rest ([10](10-tooling.md)) |
| Test framework                    | `go test` and the `testing` package, built in ([09](09-testing.md))    |
| String formatting                 | `fmt.Sprintf("%d %s %v", ...)`                                        |
| String builder                    | `strings.Builder`                                                     |
| Immutable string                  | Yes; strings are immutable byte sequences, `len` is bytes             |
| Integer division by zero          | Panics at runtime; compile error for constants                        |

## Notes by background

**Python.** Static types with inference feel like type hints that are enforced.
No list comprehensions: write the loop. No keyword arguments or defaults: use a
struct of options. Duck typing becomes implicit interfaces, which is the same
idea with compile-time checking. Virtual environments become the module cache.
`if __name__ == "__main__"` becomes `package main`.

**TypeScript / JavaScript.** No `undefined`, everything has a zero value. No
promises or event loop: blocking I/O in a goroutine is the model, and the
runtime schedules around it. Structural typing you already like; here it applies
to methods only, not fields. `npm` becomes `go get`, `node_modules` becomes a
shared cache. Objects are not maps; use a `struct` or a `map`, and choose.

**Java / C#.** Fewer keywords, no classes, no annotations, no checked
exceptions, no getters and setters by default. Interfaces are tiny and defined
by the consumer. Dependency injection is passing parameters; the repository
does it in `newContext`. No `null` safety in the type system; the nil checks are
on you. Build and test are seconds, not minutes. `Maven` / `NuGet` become
`go.mod`.

**Ruby.** Everything static and explicit; no monkey patching, no
metaprogramming. `nil` is not an object. Blocks become closures passed as
`func` values. Bundler becomes `go.mod`. The formatter is not a style guide,
it is the law.

**Rust.** No borrow checker; there is a garbage collector, and sharing is
permitted, which is why the race detector and mutexes matter. No `Result`,
`Option` or pattern matching; errors are a second return value and `nil` is the
absence. Traits become interfaces, satisfied implicitly. Fewer guarantees, much
faster compile, and a runtime that schedules goroutines for you. You will miss
`match` and enums most.

**C / C++.** Pointers without arithmetic, no manual memory management, no
header files, no macros, no undefined behaviour on integer overflow (it wraps).
Slices replace pointer-plus-length. `defer` replaces RAII. Builds are fast and
cross-compilation is a pair of environment variables.

## Habits to unlearn

1. **Wrapping everything in a class.** Functions are fine. A package of
   functions is a normal unit of design.
2. **Interfaces first.** Define the concrete type. Add an interface when a
   second implementation or a test needs one. Go's implicit satisfaction means
   you can do it later at no cost.
3. **Catching errors far from where they happen.** Handle or return each error
   at the call site. The `if err != nil` is not boilerplate; it is the design.
4. **Getters and setters.** Export the field if it is part of the API; keep it
   unexported if it is not. Methods are for behaviour.
5. **Reaching for a framework.** `net/http`, `encoding/json`, `testing` and
   `log/slog` are what most Go services run on. Add a library when you have a
   specific need it meets.
6. **Fighting the formatter.** There is nothing to configure. Set format on
   save and stop thinking about it.
7. **Deep package hierarchies.** Go packages are flat and few. `internal/` with
   a handful of focused packages is the norm, not a starting point.
8. **Inheritance to share code.** Embed a struct, or write a function that both
   call. If you are reaching for `super`, restructure.
9. **Null checks as an afterthought.** Decide which pointers can be `nil`, use
   getters on protobuf messages, return literal `nil` for interfaces.
10. **Sharing memory casually.** Anything touched by more than one goroutine
    needs a mutex, a channel, or an atomic. Run `-race` until it is a reflex.

## Idioms you will see everywhere

```go
if err != nil { return err }               // the error check
if v, ok := m[k]; ok { ... }               // comma-ok with scoped variable
defer f.Close()                            // release right after acquire
ctx context.Context                        // always the first parameter
func NewT(...) *T                          // constructor
var _ Interface = (*T)(nil)                // compile-time interface check
for _, v := range xs { ... }               // iterate values
s = append(s, v)                           // grow a slice
type X struct{ mu sync.Mutex; ... }        // mutex guards the fields below it
select { case <-ctx.Done(): return ctx.Err() ... }  // cancellation-aware wait
```

Back to the [index](README.md).
