package main

import (
	genhttpserver "github.com/vtievsky/codegen-svc/internal/services/gen-http-server"
	"github.com/vtievsky/golibs/runtime/logger"
	"go.uber.org/zap"
)

func main() {
	srv := genhttpserver.New()
	_ = srv.Stop()

	l := logger.CreateZapLogger(true, true)

	l.Info("Hello world!",
		zap.String("msg", "kjdshkjfhsdkjf"),
	)
}
