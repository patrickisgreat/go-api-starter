package main

import (
	"context"
	"fmt"
	"net"

	"github.com/soundcloud/gokit/v2/opentelemetry/instrumentation"

	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	"github.com/patrickisgreat/pb-go-api-starter/internal/config"
	"github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
	"github.com/patrickisgreat/pb-go-api-starter/internal/handlers"
)

func main() {
	ctx := context.Background()

	appserver.InitContext(ctx)

	quotesServer := api.NewQuotesServer(
		&handlers.Quotes{ServerContext: appserver.GetServerContext()},
		instrumentation.WithTwirpServerInterceptor(),
	)
	router := appserver.NewRouter(quotesServer)

	appserver.Start(ctx, &appserver.StartParams{
		Handler: instrumentation.InstrumentHandler("Quotes", router),
		Addr:    net.JoinHostPort(config.Config.Host, fmt.Sprintf("%d", config.Config.Port)),
	})
}
