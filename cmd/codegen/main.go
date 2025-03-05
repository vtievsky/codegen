package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/vtievsky/codegen-svc/gen/serverhttp"
	"github.com/vtievsky/codegen-svc/internal/config"
	"github.com/vtievsky/codegen-svc/internal/httptransport"
	specstorage "github.com/vtievsky/codegen-svc/internal/repositories/spec-storage/client"
	specstoragehttpserver "github.com/vtievsky/codegen-svc/internal/repositories/spec-storage/http-server"
	"github.com/vtievsky/codegen-svc/internal/services"
	genhttpclient "github.com/vtievsky/codegen-svc/internal/services/gen-http-client"
	genhttpserver "github.com/vtievsky/codegen-svc/internal/services/gen-http-server"
	"github.com/vtievsky/golibs/runtime/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.New()
	ctx := context.Background()
	logger := logger.CreateZapLogger(cfg.Debug, cfg.Log.EnableStacktrace)
	httpSrv := echo.New()

	serverCtx, cancel := context.WithCancel(ctx)

	specClient := specstorage.New(&specstorage.SpecStoreClientOpts{
		AccessKey: cfg.SpecStorage.AccessKey,
		SecretKey: cfg.SpecStorage.SecretKey,
		URL:       cfg.SpecStorage.URL,
	})

	httpserverSpecStore := specstoragehttpserver.New(&specstoragehttpserver.SpecHTTPServerOpts{
		Client: specClient,
	})

	genHTTPServerService := genhttpserver.New(&genhttpserver.GenHTTPServerServiceOpts{
		Logger:      logger.Named("gen-http-server"),
		SpecStorage: httpserverSpecStore,
	})

	genHTTPClientService := genhttpclient.New(&genhttpclient.GenHTTPClientServiceOpts{
		Logger:      logger.Named("gen-http-client"),
		SpecStorage: httpserverSpecStore,
	})

	services := &services.SvcLayer{
		GenHTTPServer: genHTTPServerService,
		GenHTTPClient: genHTTPClientService,
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(signalChannel) // Отмена подписки на системные события
	defer stopApp(logger, httpSrv)

	startApp(cancel, logger, httpSrv, services, cfg.Port)

	for {
		select {
		case <-signalChannel:
			logger.Info("interrupted by a signal")

			return
		case <-serverCtx.Done():
			return
		}
	}
}

func stopApp(logger *zap.Logger, httpSrv *echo.Echo) {
	defer func(alogger *zap.Logger) {
		alogger.Debug("sync zap logs")

		_ = alogger.Sync()
	}(logger)

	if err := httpSrv.Close(); err != nil {
		logger.Error("failed to close http server",
			zap.Error(err),
		)
	}
}

func startApp(
	cancel context.CancelFunc,
	logger *zap.Logger,
	httpSrv *echo.Echo,
	services *services.SvcLayer,
	port int,
) {
	defer cancel()

	serverhttp.RegisterHandlers(httpSrv, serverhttp.NewStrictHandler(
		httptransport.New(services),
		[]serverhttp.StrictMiddlewareFunc{},
	))

	address := fmt.Sprintf(":%d", port)

	if err := httpSrv.Start(address); err != nil {
		logger.Fatal("error while serve http server",
			zap.Error(err),
		)
	}
}
