package appserver

import (
	"context"
	"log"
	"sync"

	"github.com/soundcloud/gokit/v2/logger"
	"github.com/soundcloud/gokit/v2/opentelemetry/otelsetup"

	"github.com/patrickisgreat/pb-go-api-starter/internal/config"
	"github.com/patrickisgreat/pb-go-api-starter/internal/quoterepo"
	"github.com/patrickisgreat/pb-go-api-starter/internal/quoteservice"
)

type ServerContext struct {
	QuoteRepository *quoterepo.QuoteRepository
	QuoteService    *quoteservice.QuoteService
}

var sctx *ServerContext
var shutdownOtel func(context.Context) error
var initOnce sync.Once

func InitContext(ctx context.Context) {
	initOnce.Do(func() {
		var err error
		shutdownOtel, err = otelsetup.SetupSDK(ctx, otelsetup.WithECSDetector())
		if err != nil {
			log.Fatalf("Failed to set up otel: %v", err)
		}
		config.Load(ctx)
		logger.Init(logger.LoggerConfig{Config: &config.Config.Base})
		sctx = newContext()
	})
}

func GetServerContext() *ServerContext {
	if sctx == nil {
		log.Fatal("Server context was not initialized")
	}
	return sctx
}

func newContext() *ServerContext {
	quoteRepo := quoterepo.NewQuoteRepository()
	quoteService := quoteservice.NewQuoteService(quoteRepo)

	return &ServerContext{
		QuoteRepository: quoteRepo,
		QuoteService:    quoteService,
	}
}
