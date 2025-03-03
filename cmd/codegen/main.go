package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vtievsky/codegen-svc/gen/clienthttp"
	"github.com/vtievsky/codegen-svc/gen/serverhttp"
	"github.com/vtievsky/codegen-svc/internal/config"
	"github.com/vtievsky/codegen-svc/internal/httptransport"
	specstorage "github.com/vtievsky/codegen-svc/internal/repositories/spec-storage/client"
	specstoragehttpserver "github.com/vtievsky/codegen-svc/internal/repositories/spec-storage/http-server"
	"github.com/vtievsky/codegen-svc/internal/services"
	genhttpserver "github.com/vtievsky/codegen-svc/internal/services/gen-http-server"
	"github.com/vtievsky/golibs/runtime/logger"
)

const (
	specPath = "../../testdata/openapi/docs/gen-http-server.yaml"
)

func main() {
	e := echo.New()
	ctx := context.Background()
	conf := config.New()
	logger := logger.CreateZapLogger(conf.Debug, conf.Log.EnableStacktrace)

	specClient := specstorage.New(&specstorage.SpecStoreClientOpts{
		AccessKey: conf.SpecStorage.AccessKey,
		SecretKey: conf.SpecStorage.SecretKey,
		URL:       conf.SpecStorage.URL,
	})

	httpserverSpecStore := specstoragehttpserver.New(&specstoragehttpserver.SpecHTTPServerOpts{
		Client: specClient,
	})

	genHTTPServerService := genhttpserver.New(&genhttpserver.GenHTTPServerServiceOpts{
		Logger:      logger.Named("gen-http-server"),
		SpecStorage: httpserverSpecStore,
	})

	serverhttp.RegisterHandlers(e, serverhttp.NewStrictHandler(
		httptransport.New(&services.SvcLayer{
			GenHTTPServer: genHTTPServerService,
		}),
		[]serverhttp.StrictMiddlewareFunc{},
	))

	go e.Start("127.0.0.1:8080")

	time.Sleep(time.Second * 5)

	// Клиентское приложение открывает файл спецификации
	data, err := os.ReadFile("../../docs/openapi/swagger.yaml")
	if err != nil {
		log.Fatal("ошибка чтения спецификации по указанному пути")
	}

	cli, err := clienthttp.NewClientWithResponses("http://127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cli.UploadSpecHttpWithResponse(ctx, clienthttp.UploadSpecHttpRequest{
		Name: "codegen",
		Spec: data,
	})
	if err != nil {
		log.Fatal(err, resp.Status())
	}

	respCli, err := cli.GenerateSpecServerHttpWithResponse(ctx, &clienthttp.GenerateSpecServerHttpParams{
		Name: "codegen",
	})
	if err != nil {
		log.Fatal(err)
	}

	outputDir := fmt.Sprintf("./tmp/%s", "codegen")

	// Клиентское приложение удаляет предыдущую версию файла спецификации
	if err := os.RemoveAll(outputDir); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Клиентское приложение сохраняет новую версию файла спецификации
	outputFile := path.Join(outputDir, "clienthttp.go")

	err = os.WriteFile(outputFile, respCli.JSON200.Spec, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("exit")
}
