package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/vtievsky/codegen-svc/internal/config"
	genhttpserver "github.com/vtievsky/codegen-svc/internal/services/gen-http-server"
	"github.com/vtievsky/golibs/runtime/logger"
)

const (
	specPath  = "../../testdata/openapi/docs/gen-http-server.yaml"
	outputDir = "../../gen/serverhttp"
)

func main() {
	cfg := config.New()
	ctx := context.Background()
	logger := logger.CreateZapLogger(cfg.Debug, cfg.Log.EnableStacktrace)

	httpserver := genhttpserver.New(&genhttpserver.GenHTTPServerServiceOpts{
		Logger: logger.Named("gen-http-server"),
	})

	// Клиентское приложение открывает файл спецификации
	data, err := os.ReadFile(specPath)
	if err != nil {
		log.Fatal("ошибка чтения спецификации по указанному пути")
	}

	// Клиентское приложение "отправляет" файл спецификации серверу генерации
	resp, err := httpserver.GenHTTPServer(ctx, data)
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
