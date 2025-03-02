package main

import (
	"time"

	"github.com/vtievsky/codegen-svc/internal/config"
	genhttpserver "github.com/vtievsky/codegen-svc/internal/services/gen-http-server"
	"github.com/vtievsky/golibs/runtime/logger"
)

func main() {
	cfg := config.New()
	logger := logger.CreateZapLogger(cfg.Debug, cfg.Log.EnableStacktrace)

	srv := genhttpserver.New(&genhttpserver.GenHTTPServerServiceOpts{
		Logger: logger,
	})

	_ = srv.Start()
	time.Sleep(time.Second * 3)
	srv.Stop()
}
