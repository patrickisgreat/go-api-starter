# API guide

The service exposes one Twirp service, `Quotes`. This page explains how to call it
and how to extend it. The field-by-field reference generated from the proto files
is in [pb_go_api_starter.md](pb_go_api_starter.md).

## Routing

Twirp maps every RPC to:

```
POST /twirp/<package>.<Service>/<Method>
```

For this service the package is `proto.patrickisgreat.pb_go_api_starter.api`, so
the three routes are:

```
POST /twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuote
POST /twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuoteFlaky
POST /twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuoteLate
```

Send `Content-Type: application/json` for JSON bodies or
`application/protobuf` for binary. JSON field names are the proto field names in
either `snake_case` or `lowerCamelCase`; responses use `lowerCamelCase` unless the
client asks otherwise.

## RPCs

### GetQuote

Returns a random quote.

```sh
curl -s -X POST -H "Content-Type: application/json" \
  http://localhost:8000/twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuote \
  --data '{"userSession": {"userUrn": "myapp:users:42"}}'
```

```json
{"quote": "Simplicity is prerequisite for reliability.\n        ― Edsger W. Dijkstra"}
```

### GetQuoteFlaky

Returns a random quote but fails roughly half the time with an `internal` error.
Set `fail` to `true` to make it fail every time. It exists so that error-rate
alerts and dashboards can be exercised deliberately.

```sh
curl -s -X POST -H "Content-Type: application/json" \
  http://localhost:8000/twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuoteFlaky \
  --data '{"fail": true}'
```

```json
{"code": "internal", "msg": "flaky error"}
```

### GetQuoteLate

Returns a random quote after `delayMs` milliseconds. When `delayMs` is omitted the
delay is random between 100 and 1000 ms. It exists to exercise latency alerts.

```sh
curl -s -X POST -H "Content-Type: application/json" \
  http://localhost:8000/twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuoteLate \
  --data '{"delayMs": 2000}'
```

## User sessions

Every request message has an optional `userSession` of type `UserSession`
(defined in `rpc/common_session.proto`). It carries the caller's identity and
context: user URN, agent URN, geo, OAuth scopes and feature flags.

Only one rule is enforced today: when `userSession.userUrn` is set, it must parse
as a URN whose collection is `users`. Anything else is rejected:

```json
{"code": "invalid_argument", "msg": "userSession.userUrn invalid user urn", "meta": {"argument": "userSession.userUrn"}}
```

## Errors

Twirp errors are JSON with `code`, `msg` and optional `meta`. The HTTP status is
derived from the code. Codes used by this service:

| Code               | HTTP | When                                             |
| ------------------ | ---- | ------------------------------------------------ |
| `invalid_argument` | 400  | User URN present but not a `users` URN           |
| `internal`         | 500  | `GetQuoteFlaky` failing, or no quotes available  |
| `bad_route`        | 404  | Unknown method, wrong HTTP verb, bad content type |

## Health check

`GET /-/health` returns `200 OK` with the body `OK`. It does not check
dependencies because there are none.

## Adding an RPC

1. **Define it in the proto.** Add request and response messages and the `rpc`
   line to the `Quotes` service in `rpc/pb_go_api_starter.proto`. Follow the
   naming rules at the top of that file: verb-first method names, and
   `<Method>Request` / `<Method>Response` message names.
2. **Regenerate.** `task gen:rpc`. The build now fails because `handlers.Quotes`
   no longer satisfies the generated interface. That is the compiler giving you a
   to-do list.
3. **Add the handler.** Create `internal/api/<method>.go` with a
   `<Method>Handler(ctx, serverContext, request)` function, and a one-line method
   on `handlers.Quotes` that delegates to it.
4. **Extend the service or repository** if the RPC needs new behaviour.
5. **Test.** A unit test in `internal/api` with a hand-built `ServerContext`, and
   an end to end test in `end_to_end/` that hits the real route.
6. **Regenerate the docs.** `task gen:docs` rewrites `pb_go_api_starter.md`.

The [Go learning path](../golang-learning/README.md) walks through exactly this
exercise step by step.
