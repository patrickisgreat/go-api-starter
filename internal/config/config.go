package config

import (
	"context"
	"log/slog"

	"github.com/soundcloud/gokit/v2/baseconfig"
	"github.com/soundcloud/gokit/v2/logger"
)

type config struct {
	baseconfig.Base

	LogPretty bool `envconfig:"LOG_PRETTY" default:"false"`
}

var Config config

func Load(ctx context.Context) {
	err := baseconfig.Load[config](&Config)
	if err != nil {
		logger.Get(ctx).Error("failed loading app config", slog.Any("err", err))
		panic("failed loading app config")
	}
}
