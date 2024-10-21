package main

import (
	"context"
	"fmt"

	"github.com/slipneff/gogger"
	"github.com/slipneff/gogger/log"
	"github.com/slipneff/auto-yt/internal/di"
	"github.com/slipneff/auto-yt/internal/utils/config"
	"github.com/slipneff/auto-yt/internal/utils/flags"
)

func main() {
	flags := flags.MustParseFlags()
	cfg := config.MustLoadConfig(flags.EnvMode, flags.ConfigPath)
	gogger.ConfigureZeroLogger()

	container := di.New(context.Background(), cfg)
	log.Info(fmt.Sprintf("Server starting at %s:%d", cfg.Host, cfg.Port))
	err := container.GetHttpServer().ListenAndServe()
	if err != nil {
		log.Panic(err, "Fail serve HTTP")
	}
}
