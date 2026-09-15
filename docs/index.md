# pb-go-api-starter

A small, production-shaped Go API that serves random quotes.

The business logic is intentionally trivial. What this repository demonstrates is
everything around it: a protobuf-defined [Twirp](https://github.com/twitchtv/twirp)
API, a layered package structure, unit and end to end tests, hot reload for local
development, a distroless container image, OpenTelemetry instrumentation and
Terraform for AWS ECS.

## Where to go

| If you want to...                                  | Read                                              |
| -------------------------------------------------- | ------------------------------------------------- |
| Run it on your machine in five minutes             | [Getting started](getting-started.md)             |
| Understand how a request moves through the code    | [Architecture](architecture.md)                   |
| Call the API or add a new RPC                      | [API guide](api.md), [proto reference](pb_go_api_starter.md) |
| Know which environment variables matter            | [Configuration](configuration.md)                 |
| Write or run tests                                 | [Testing](testing.md)                             |
| Ship it                                            | [Deployment](deployment.md)                       |
| Learn Go using this codebase                       | [Go learning path](../golang-learning/README.md)  |

## The service in one paragraph

`cmd/app/main.go` loads configuration from environment variables, initialises
OpenTelemetry and logging, builds a `ServerContext` holding a quote repository and
quote service, wraps the generated Twirp server in instrumentation, and starts an
HTTP server that shuts down gracefully on `SIGINT` or `SIGTERM`. The `Quotes`
service has three RPCs: one that works, one that fails on purpose about half the
time, and one that answers slowly on purpose. The latter two exist so that
alerting and dashboards can be exercised without a real incident.
