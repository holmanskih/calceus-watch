package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/holmanskih/calceus-watch/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	projectPath = "/Users/holmanskih/Desktop/calceus/calceus-watch/test_data/"
	buildPath   = "/Users/holmanskih/Desktop/calceus/calceus-watch/test_data/build/"
)

func main() {
	cfg := internal.Config{
		ProjectDir: projectPath,
		BuildDir:   buildPath,
		Mode:       internal.ModeProduction,
	}

	log, err := initLogger()
	if err != nil {
		panic(fmt.Sprintf("init logger err %e", err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//comiler spool
	pool := internal.NewCompilerPool(log)
	go pool.Run(ctx, cfg)

	// out compiler bus to communicate parser and compiler files
	compilerChan := make(chan internal.Compiler)

	// start parser worker
	parser := internal.NewParser(cfg, log)
	go parser.Watch(ctx, cancel, compilerChan, pool.GetNewCompilerOutBus())

	// start compilerChan worker pool

	for {
		select {
		case <-ctx.Done():
			log.Info("calceus watch was gracefully shutdown")
			return

		case <-done:
			cancel()
		}
	}
}

func initLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC822)
	return cfg.Build()
}
