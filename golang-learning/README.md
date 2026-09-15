# Learning Go with pb-go-api-starter

A guided path through Go for developers who already know how to build software
in another language and want to become productive in Go quickly. Instead of toy
examples, every chapter points at real code in this repository, and the path
ends with you shipping a new RPC end to end.

## Who this is for

You have written production code in something like Python, TypeScript, Java, C#,
Ruby or Rust. You know what a struct, an interface, a thread and an HTTP server
are. You do not need those explained. You need to know how Go does them, where Go
differs from what you expect, and which habits to drop.

## How to use it

1. Get the service running first: [docs/getting-started.md](../docs/getting-started.md).
2. Work through the chapters in order. Each one is 15 to 40 minutes and ends
   with a short **Try it** section. Do the exercises; reading Go is not the same
   as writing Go.
3. Keep the [Go specification](https://go.dev/ref/spec) and
   [Effective Go](https://go.dev/doc/effective_go) open in another tab. They are
   short, and this path tells you which sections matter when.
4. Finish with the capstone, which adds a real RPC to the service with tests and
   docs.

## Chapters

| #  | Chapter                                                       | What you will be able to do                                       |
| -- | ------------------------------------------------------------- | ----------------------------------------------------------------- |
| 01 | [Toolchain and modules](01-toolchain-and-modules.md)          | Build, test and run Go code; understand `go.mod` and `go.sum`     |
| 02 | [Packages and project layout](02-packages-and-layout.md)      | Read `cmd/` and `internal/`; know what is exported and why        |
| 03 | [Syntax essentials](03-syntax-essentials.md)                  | Read any file in this repo without stopping                       |
| 04 | [Errors](04-errors.md)                                        | Return, wrap, inspect and translate errors the Go way             |
| 05 | [Structs, methods and interfaces](05-structs-methods-interfaces.md) | Model the service's types; satisfy an interface without declaring it |
| 06 | [Context, goroutines and shutdown](06-context-goroutines-shutdown.md) | Follow `ctx` through the code; understand the server lifecycle |
| 07 | [HTTP, protobuf and Twirp](07-http-protobuf-twirp.md)         | Know what the generated code does and what is yours to write      |
| 08 | [Standard library tour](08-standard-library-tour.md)          | Use `embed`, `encoding/json`, `log/slog`, `math/rand` and friends  |
| 09 | [Testing](09-testing.md)                                      | Write unit and end to end tests in the repo's style                |
| 10 | [Tooling and code quality](10-tooling.md)                     | Use `gofmt`, `goimports`, `go vet`, `golangci-lint` and `air`      |
| 11 | [Capstone: add an RPC](11-capstone-add-an-rpc.md)             | Ship `ListQuotes` from proto to end to end test                    |
| 12 | [Further exercises](12-further-exercises.md)                  | Keep going: persistence, middleware, concurrency, benchmarks       |
| —  | [Coming from another language](coming-from-another-language.md) | A cheat sheet of Go equivalents and the habits to unlearn        |

## The mental model in five lines

- Go is small. The whole language fits in your head; the standard library is where
  the depth is.
- Errors are values that you return and check. There are no exceptions.
- Interfaces are satisfied implicitly. You never write `implements`.
- Concurrency is goroutines plus channels plus `context.Context` for cancellation.
- Formatting, vetting and dependency management are built in and not up for
  debate. That is a feature.
