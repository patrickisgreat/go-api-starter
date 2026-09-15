# Getting started

This page gets the service running locally and confirms it answers requests.

## 1. Install prerequisites

- [Go](https://go.dev/doc/install) 1.23 or newer. Check with `go version`.
- [Task](https://taskfile.dev), the command runner used for every workflow here:

  ```sh
  go install github.com/go-task/task/v3/cmd/task@latest
  ```

  Make sure `$(go env GOPATH)/bin` is on your `PATH` so that `task` and the
  `protoc` plugins installed below can be found.

- Docker, only for building the container image or regenerating the API docs.

## 2. Private module access

The service depends on `github.com/soundcloud/gokit/v2`, which lives in a private
GitHub repository. Go must be told not to use the public proxy for it, and `git`
must be able to authenticate:

```sh
export GOPRIVATE=github.com/soundcloud/*
# use SSH for GitHub so that go mod download can clone the private module
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

If you do not have access, `go build` fails with `Repository not found`. That is
the symptom to look for.

## 3. Install the toolchain

```sh
task deps:macos
```

On macOS this installs `protobuf` and `golangci-lint` with Homebrew, then
`air`, `goimports`, `protoc-gen-go` and `protoc-gen-twirp` with `go install`.
On Linux install `protoc` and `golangci-lint` from your package manager and run
`task deps:proto` for the Go tools.

## 4. Generate code and download dependencies

```sh
task gen:rpc    # protoc -> internal/generated/pb_go_api_starter/api
task install    # go get -t ./...
```

`internal/generated` is git ignored. Every fresh clone needs `task gen:rpc`
before anything compiles; the build task runs it for you.

## 5. Run it

```sh
task run
```

This runs `air`, which builds the binary, starts it, and rebuilds on every save.
Configuration comes from `.env/.env.local`: the server listens on
`localhost:8000` with pretty debug logs and telemetry export disabled.

## 6. Make a request

```sh
curl -i http://localhost:8000/-/health

curl -s -X POST -H "Content-Type: application/json" \
  http://localhost:8000/twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuote \
  --data '{}'
```

You should see `OK` from the health check and a JSON object with a `quote` field
from the RPC.

## 7. Run the tests

```sh
task test        # unit tests, preceded by tidy, lint and vet
task test:e2e    # builds the binary, starts it, runs end_to_end/
```

## Common problems

| Symptom                                              | Fix                                                        |
| ---------------------------------------------------- | ---------------------------------------------------------- |
| `package .../internal/generated/... is not in std`   | Run `task gen:rpc`.                                        |
| `protoc-gen-go: program not found`                   | Add `$(go env GOPATH)/bin` to `PATH`, or rerun `task deps:proto`. |
| `Repository not found` for `soundcloud/gokit`        | See [private module access](#2-private-module-access).     |
| `task: command not found`                            | Same `PATH` issue as above.                                |
| Port 8000 already in use                             | Change `PORT` in `.env/.env.local`.                        |
