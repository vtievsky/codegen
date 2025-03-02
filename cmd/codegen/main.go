package main

import (
	"os"

	genhttpserver "github.com/vtievsky/codegen-svc/internal/services/gen-http-server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateZapLogger(debug, stacktrace bool) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	stacktraceLevel := zapcore.ErrorLevel
	if !stacktrace {
		stacktraceLevel = zapcore.FatalLevel + 1
	}

	stdoutFiler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		if debug {
			return level < zapcore.ErrorLevel
		}

		return level > zapcore.DebugLevel && level < zapcore.ErrorLevel
	})

	stderrFilter := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			stdoutFiler,
		),
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stderr),
			stderrFilter,
		),
	)

	return zap.New(core, zap.AddStacktrace(stacktraceLevel))
}

func main() {
	srv := genhttpserver.New()
	_ = srv.Stop()

	l := CreateZapLogger(true, true)

	l.Info("Hello world!",
		zap.String("msg", "kjdshkjfhsdkjf"),
	)
}
