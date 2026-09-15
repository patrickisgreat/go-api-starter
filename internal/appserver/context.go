package appserver

import (
	"context"
	"log"
	"sync"

	"github.com/soundcloud/gokit/v2/logger"
	"github.com/soundcloud/gokit/v2/opentelemetry/otelsetup"

	"github.com/patrickisgreat/pb-go-api-starter/internal/config"
	"github.com/patrickisgreat/pb-go-api-starter/internal/cookierepo"
	"github.com/patrickisgreat/pb-go-api-starter/internal/cookieservice"
)

type ServerContext struct {
	CookieRepository *cookierepo.CookieRepository
	CookieService    *cookieservice.CookieService
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
	cookieRepo := cookierepo.NewCookieRepository()
	cookieService := cookieservice.NewCookieService(cookieRepo)

	return &ServerContext{
		CookieRepository: cookieRepo,
		CookieService:    cookieService,
	}
}
