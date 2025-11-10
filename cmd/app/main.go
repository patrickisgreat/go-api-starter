package main

import (
	"context"
	"fmt"
	"net"

	"github.com/soundcloud/gokit/v2/opentelemetry/instrumentation"

	"github.com/patrickisgreat/go-api-starter/internal/appserver"
	"github.com/patrickisgreat/go-api-starter/internal/config"
	"github.com/patrickisgreat/go-api-starter/internal/generated/go_api_starter/api"
	"github.com/patrickisgreat/go-api-starter/internal/handlers"
)

func main() {
	ctx := context.Background()

	appserver.InitContext(ctx)

	fortuneServer := api.NewFortuneServer(
		&handlers.Fortune{ServerContext: appserver.GetServerContext()},
		instrumentation.WithTwirpServerInterceptor(),
	)
	router := appserver.NewRouter(fortuneServer)

	appserver.Start(ctx, &appserver.StartParams{
		Handler: instrumentation.InstrumentHandler("Fortune", router),
		Addr:    net.JoinHostPort(config.Config.Host, fmt.Sprintf("%d", config.Config.Port)),
	})
}
