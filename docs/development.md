# Development workflow

## Everyday commands

Every workflow is a [Task](https://taskfile.dev) task defined in `Taskfile.yml`.
`task --list-all` shows them all.

| Task                | Runs                                                        | Notes                                              |
| ------------------- | ----------------------------------------------------------- | -------------------------------------------------- |
| `task run`          | `air`                                                       | Hot reload. Rebuilds and restarts on save          |
| `task build`        | `task gen:rpc` then `go build -o cmd/app/bin/pb-go-api-starter` | Static binary, `GOOS` defaults to `darwin`     |
| `task gen:rpc`      | `protoc` with the Go and Twirp plugins                      | Wipes and recreates `internal/generated`           |
| `task gen:docs`     | `protoc-gen-doc` in Docker                                  | Rewrites `docs/pb_go_api_starter.md`               |
| `task format`       | `goimports -w -local <module>`                              | Groups imports: stdlib, third party, this module   |
| `task lint`         | `golangci-lint run`                                         | Runs `format` first                                |
| `task vet`          | `go vet ./...`                                              |                                                    |
| `task tidy`         | `go mod tidy`                                               |                                                    |
| `task test`         | `tidy`, `lint`, `vet`, then `go test`                       | Unit tests only                                    |
| `task test:e2e`     | `build`, then `go test ./end_to_end`                        | Spawns the binary                                  |
| `task upgrade`      | `go get -u ./...`                                           | Follow with `task tidy` and `task test`            |

Tasks declare `sources` and `generates`, so Task skips work whose inputs have not
changed. Delete `.task/` to force a full rerun.

## Hot reload

`.air.toml` tells `air` to watch `app`, `internal`, `pkg` and `rpc` for `.go`
changes, run `task build`, and restart `cmd/app/bin/pb-go-api-starter`. Test
files and `internal/generated` are excluded from the watch. Build errors go to
`build-errors.log`.

## Changing the API

1. Edit `rpc/pb_go_api_starter.proto`.
2. `task gen:rpc`.
3. Fix whatever the compiler now complains about in `internal/handlers`.
4. Implement, test, and run `task gen:docs`.

The [API guide](api.md#adding-an-rpc) has the detailed version.

## Code style

- `gofmt` formatting is not negotiable; `goimports` enforces it plus import
  grouping. Run `task format` before committing.
- Packages are named for what they provide (`quoterepo`, `quoteservice`), all
  lower case, no underscores.
- Handlers take `ctx` first and return `(response, error)`.
- Errors are returned, not logged and swallowed. Log at the point where the error
  is handled, with `logger.Get(ctx)` so that log lines carry trace context.
- Anything under `internal/` cannot be imported from outside this module. That is
  a Go rule, not a convention, and it is why almost everything lives there.

## Dependencies

`go.mod` pins direct dependencies; `go.sum` pins checksums for everything. After
adding an import, run `task tidy` so both files are updated. The
`github.com/soundcloud/gokit/v2` module is private; see
[Getting started](getting-started.md#2-private-module-access).

## Backstage

`catalog-info.yaml` registers the service and its API with Backstage, and
`mkdocs.yml` configures this documentation site for TechDocs. When you add a
page to `docs/`, add it to the `nav` in `mkdocs.yml`.
