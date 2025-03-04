package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type LogConfig struct {
	EnableStacktrace bool `envconfig:"CODEGEN_LOG_ENABLE_STACKTRACE" default:"false"`
}

type SpecStorageConfig struct {
	URL       string `envconfig:"CODEGEN_SPECSTORAGE_URL" required:"true"`
	AccessKey string `envconfig:"CODEGEN_SPECSTORAGE_ACCESS_KEY" required:"true"`
	SecretKey string `envconfig:"CODEGEN_SPECSTORAGE_SECRET_KEY" required:"true"`
}

type Config struct {
	Debug bool `envconfig:"CODEGEN_DEBUG" default:"false"`

	Log         LogConfig
	SpecStorage SpecStorageConfig
}

func New() *Config {
	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		err = fmt.Errorf("error while parse env config | %w", err)

		log.Fatal(err)
	}

	return cfg
}
