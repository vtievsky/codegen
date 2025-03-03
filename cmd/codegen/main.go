package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/vtievsky/codegen-svc/internal/config"
	specstorage "github.com/vtievsky/codegen-svc/internal/repositories/spec-storage"
	specstoragehttpserver "github.com/vtievsky/codegen-svc/internal/repositories/spec-storage/http-server"
	genhttpserver "github.com/vtievsky/codegen-svc/internal/services/gen-http-server"
	"github.com/vtievsky/golibs/runtime/logger"
)

const (
	specPath  = "../../testdata/openapi/docs/gen-http-server.yaml"
	outputDir = "../../gen/serverhttp"
)

func main() {
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

	httpserver := genhttpserver.New(&genhttpserver.GenHTTPServerServiceOpts{
		Logger:      logger.Named("gen-http-server"),
		SpecStorage: httpserverSpecStore,
	})

	// // Клиентское приложение открывает файл спецификации
	// data, err := os.ReadFile(specPath)
	// if err != nil {
	// 	log.Fatal("ошибка чтения спецификации по указанному пути")
	// }

	// // Клиентское приложение "отправляет" файл спецификации серверу генерации
	// err = httpserver.SaveSpec(ctx, "auth-id", data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Клиентское приложение "отправляет" файл спецификации серверу генерации
	resp, err := httpserver.GenerateCode(ctx, "auth-id")
	if err != nil {
		log.Fatal(err)
	}

	// Клиентское приложение удаляет предыдущую версию файла спецификации
	if err := os.RemoveAll(outputDir); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Клиентское приложение сохраняет новую версию файла спецификации
	outputFile := path.Join(outputDir, "serverhttp.go")

	err = os.WriteFile(outputFile, resp, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("exit")
}
