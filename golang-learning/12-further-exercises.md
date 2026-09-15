# 12. Further exercises

Each exercise extends the service in a direction real services grow. They are
ordered roughly by effort. Pick the ones that match what you need to learn; you
do not have to do them in sequence. For each, write tests, run `task test`, and
commit in small steps.

## 1. Doc comments everywhere

Add doc comments to every exported identifier under `internal/`, enable the
`revive` linter, and make `task lint` pass. Then run `go doc ./internal/...`
and read your own documentation.

*Practises:* doc conventions, linter configuration.

## 2. `GetQuoteByIndex`

An RPC that returns the quote at a given position, with `not_found` for out of
range and `invalid_argument` for negative. Include the index and the total
count in the response.

*Practises:* Twirp error codes, bounds checking, the compile-driven workflow
from the capstone.

## 3. Deterministic randomness

`GetQuote` and `GetQuoteFlaky` use the global random source, which makes their
random paths untestable. Add a `Rand *rand.Rand` field to `ServerContext`
(or an interface with `Intn`), thread it through, and write tests that assert
the exact quote chosen with a seeded source.

*Practises:* dependency injection without a framework, interfaces where they
earn their place, `math/rand`.

## 4. Configuration for behaviour

Make the flaky failure rate and the late delay range configurable through
environment variables (`FLAKY_FAILURE_RATE`, `LATE_MIN_MS`, `LATE_MAX_MS`) with
sensible defaults, following the `LogPretty` pattern in `internal/config`.
Document them in `docs/configuration.md`.

*Practises:* struct tags, `envconfig`, keeping configuration in one place.

## 5. An in-memory write path

Add `AddQuote` (append a quote at runtime) and make `ListQuotes` see it.
Handlers run concurrently, so `QuoteRepository` now needs a `sync.RWMutex`.
Prove it with a test that adds and reads from several goroutines under
`go test -race`.

*Practises:* mutexes, the race detector, the difference between read and write
locks.

## 6. Middleware

Write `func RequestLogger(next http.Handler) http.Handler` that logs method,
path, status and duration for every request using `slog`. Capturing the status
requires wrapping `http.ResponseWriter`; that is the interesting part. Wire it
in `main.go` and then move it into `appserver`.

*Practises:* the `Handler` interface, wrapping types to satisfy interfaces,
closures.

## 7. Timeouts and cancellation

Give every RPC a deadline via a Twirp server hook or middleware that derives
`context.WithTimeout`. Make `GetQuoteLate` respect cancellation by replacing
`time.Sleep` with a `select` on `time.After` and `ctx.Done()`, returning
`twirp.DeadlineExceeded`. Test with a short deadline.

*Practises:* `context`, `select`, cooperative cancellation.

## 8. A Twirp client binary

Create `cmd/quotectl/main.go`: a command line client using
`api.NewQuotesJSONClient` with subcommands `get`, `list`, `flaky` and `late`,
parsed with the standard `flag` package. Print typed errors sensibly.

*Practises:* generated clients, `flag`, `os.Args`, multiple `main` packages in
one module.

## 9. Persistence

Replace the embedded JSON with SQLite via `database/sql` and
`modernc.org/sqlite` (pure Go, keeps `CGO_ENABLED=0`). Keep the
`QuoteRepository` method signatures unchanged so that nothing above it moves.
Load the JSON into the database at start-up if it is empty.

*Practises:* `database/sql`, `context`-aware queries, the payoff of a
repository layer.

## 10. Graceful shutdown under load

Start the server, fire a stream of `GetQuoteLate` requests with a load tool
(`hey` or `wrk`), send `SIGTERM`, and confirm from the logs that in-flight
requests complete and none fail. Then set the shutdown budget to one second and
watch what changes.

*Practises:* the server lifecycle, reading logs, understanding what
`Shutdown` promises.

## 11. Benchmarks and profiling

Benchmark `QuoteRepository.Find` from the capstone. Then add
`import _ "net/http/pprof"` behind a config flag, run a load test, and capture
a CPU profile with `go tool pprof`. Find the hot spot in `Find` (hint:
`strings.ToLower` on every quote, every call) and precompute it.

*Practises:* `testing.B`, `pprof`, measuring before optimising.

## 12. Generics, once

Write `func Map[T, U any](s []T, f func(T) U) []U` and `func Filter[T any](s []T, keep func(T) bool) []T`
in a small `internal/sliceutil` package, with tests, and use `Filter` in
`Find`. Then decide whether it made the code clearer. Either answer is fine;
knowing when generics help is the skill.

*Practises:* type parameters, constraints, judgement.

## 13. Multi-stage CI

Write a GitHub Actions workflow that runs the CI sequence from
[chapter 10](10-tooling.md) on every pull request, caches the module download,
and builds the container image on merge to `main`.

*Practises:* the full toolchain in an unattended environment, build caching,
`GOPRIVATE` in CI.

## Where to go after this

- [The Go Programming Language](https://www.gopl.io) (Donovan and Kernighan)
  remains the best thorough book.
- [100 Go Mistakes and How to Avoid Them](https://100go.co) is the best second
  book: every entry is a thing you will otherwise learn the hard way.
- [Go by Example](https://gobyexample.com) for quick lookups.
- Read the standard library source. `net/http/server.go` and `sync/once.go`
  are short, well commented, and written by people who know the language.
