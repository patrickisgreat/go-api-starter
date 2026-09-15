# 11. Capstone: add an RPC

You will add `ListQuotes` to the `Quotes` service: an RPC that returns up to
`limit` quotes, optionally filtered to those containing a search term. It
touches every layer, needs new tests, and exercises everything from the earlier
chapters. Budget one to two hours. Do it on a branch and open a pull request at
the end, even if only for yourself.

## Target behaviour

```
POST /twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/ListQuotes
{"limit": 3, "contains": "computer"}

200 {"quotes": ["...", "...", "..."], "total": 41}
```

- `limit` defaults to 10 when omitted and is capped at 100. A negative value is
  an `invalid_argument` error.
- `contains` is a case-insensitive substring filter; empty means no filter.
- `total` is the number of quotes that matched before the limit was applied.
- The same `user_session` validation as the other RPCs.

## Step 1: define the contract

In `rpc/pb_go_api_starter.proto` add the messages and the RPC. Follow the naming
convention in the comment at the top of the file.

```proto
message ListQuotesRequest {
    // The session of the user requesting quotes.
    proto.patrickisgreat.common.session.UserSession user_session = 1;
    /* Maximum number of quotes to return. Defaults to 10, capped at 100. */
    google.protobuf.Int32Value limit = 2;
    /* Case-insensitive substring filter. Empty means no filter. */
    string contains = 3;
}

message ListQuotesResponse {
    /* Matching quotes, at most `limit` of them */
    repeated string quotes = 1;
    /* Number of quotes that matched before the limit was applied */
    int32 total = 2;
}
```

and inside `service Quotes`:

```proto
    /* List quotes, optionally filtered by substring. */
    rpc ListQuotes (ListQuotesRequest) returns (ListQuotesResponse) {
    }
```

Why `Int32Value` for `limit` but plain `string` for `contains`? A missing
`contains` and an empty `contains` mean the same thing, so the zero value is
fine. A missing `limit` should mean "default", which is different from `0`, so
it needs a type that can be `nil`. [Chapter 07](07-http-protobuf-twirp.md)
covers this.

## Step 2: regenerate and let the compiler guide you

```sh
task gen:rpc
go build ./...
```

The build fails in `cmd/app/main.go`: `*handlers.Quotes does not implement
api.Quotes (missing method ListQuotes)`. That message is your to-do list.

## Step 3: the repository

Add to `internal/quoterepo/repository.go`:

```go
// Find returns quotes containing substr (case-insensitive; empty matches all)
// up to limit, and the total number that matched.
func (cr *QuoteRepository) Find(ctx context.Context, substr string, limit int) ([]string, int) {
	needle := strings.ToLower(substr)
	matched := make([]string, 0, limit)
	total := 0
	for _, q := range cr.Quotes {
		if needle != "" && !strings.Contains(strings.ToLower(q), needle) {
			continue
		}
		total++
		if len(matched) < limit {
			matched = append(matched, q)
		}
	}
	return matched, total
}
```

Things to notice: `make([]string, 0, limit)` pre-sizes the slice; `continue`
keeps the happy path unindented; the function returns two values and no error
because nothing here can fail. Write the test first if you prefer:
`internal/quoterepo/repository_test.go` with a table covering empty filter,
case-insensitivity, limit smaller than matches, and no matches.

## Step 4: the service

Add to `internal/quoteservice/service.go`:

```go
const (
	DefaultLimit = 10
	MaxLimit     = 100
)

func (cs *QuoteService) ListQuotes(ctx context.Context, contains string, limit int) ([]string, int, twirp.Error) {
	if limit < 0 {
		return nil, 0, twirp.InvalidArgumentError("limit", "must not be negative")
	}
	if limit == 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	quotes, total := cs.repo.Find(ctx, contains, limit)
	return quotes, total, nil
}
```

Defaults and caps are business rules, so they live here rather than in the
handler or the repository. The service returns `twirp.Error` like its
neighbours. Add a test for each branch of the limit logic.

## Step 5: the handler

Create `internal/api/list.go`:

```go
package api

import (
	"context"

	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

func ListQuotesHandler(ctx context.Context, serverContext *appserver.ServerContext, request *proto.ListQuotesRequest) (*proto.ListQuotesResponse, error) {
	if err := validateUserSession(request.GetUserSession()); err != nil {
		return nil, err
	}

	limit := int(request.GetLimit().GetValue())   // 0 when limit was not sent

	quotes, total, err := serverContext.QuoteService.ListQuotes(ctx, request.GetContains(), limit)
	if err != nil {
		return nil, err
	}

	return &proto.ListQuotesResponse{
		Quotes: quotes,
		Total:  int32(total),
	}, nil
}
```

Then the one-liner in `internal/handlers/handlers.go`:

```go
func (h *Quotes) ListQuotes(ctx context.Context, request *proto.ListQuotesRequest) (*proto.ListQuotesResponse, error) {
	return api.ListQuotesHandler(ctx, h.ServerContext, request)
}
```

`go build ./...` should now pass. Pause on the `err != nil` after
`ListQuotes`: `err` is a `twirp.Error`, and returning it through a function
whose result type is `error` is exactly the nil interface trap from
[chapter 04](04-errors.md). The code above is safe because it returns `nil`
literally on the success path. Had it been `return resp, err` with a nil
`twirp.Error`, every successful call would have become a 500. Try it and see.

## Step 6: unit tests

`internal/api/list_test.go`, following `quote_test.go`:

- Happy path: a repository with three quotes, `contains` matching two, assert
  the two come back and `total` is 2.
- Limit: five quotes, `limit: 2`, assert two quotes and `total` 5.
- Negative limit: assert the error string starts with
  `twirp error invalid_argument: limit`.
- Invalid session: copy the existing pattern.

Run `task test`. Fix lint findings; `errcheck` and `goimports` are the usual
ones.

## Step 7: end to end test

Add to `end_to_end/quote_test.go`:

```go
func TestListQuotes(t *testing.T) {
	c := qt.New(t)

	url := apiHost + "/twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/ListQuotes"
	payload := map[string]interface{}{
		"limit":    2,
		"contains": "the",
	}
	response, status, err := postRequest(url, payload)
	c.Assert(err, qt.IsNil)
	c.Assert(status, qt.Equals, 200)
	quotes, ok := response["quotes"].([]interface{})
	c.Assert(ok, qt.IsTrue)
	c.Assert(quotes, qt.HasLen, 2)
	c.Assert(response["total"].(float64) > 2, qt.IsTrue)
}
```

Note the `float64`: JSON numbers decode to `float64` in an untyped map.
Run `task test:e2e`.

## Step 8: docs and PR

```sh
task gen:docs
```

regenerates `docs/pb_go_api_starter.md`. Add a `ListQuotes` section to
`docs/api.md` with a curl example, and a row to the RPC table in the README.
Commit in small steps (proto and generated interface; repository; service;
handler; tests; docs) so that each commit is reviewable on its own.

## What you practised

- Protobuf messages and the difference between zero values and absence
- Code generation and using the compiler as a checklist
- Slices, `range`, `strings`, pre-sizing with `make`
- Business rules in the service layer, error translation to Twirp codes
- The nil interface trap, avoided deliberately
- Unit tests with a hand-built `ServerContext`, table tests, end to end tests
- The lint and format loop

## Read next

- [12. Further exercises](12-further-exercises.md)
