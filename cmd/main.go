package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/holmanskih/calceus-watch/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogLevelInfo  = "info"
	LogLevelDebug = "debug"
)

func main() {
	projectPath := flag.String("projectPath", "", "project absolute path")
	buildPath := flag.String("buildPath", "", "project build absolute path")

	// optional
	sassDirPath := flag.String("sassDirPath", "scss", "sass directory relative path")
	mode := flag.Bool("prod", false, "production run mode")
	logLevel := flag.String("log", LogLevelDebug, "log level")

	flag.Parse()

	var runMode internal.Mode
	if *mode {
		runMode = internal.ModeProduction
	} else {
		runMode = internal.ModeDevelopment
	}

	log, err := initLogger(*logLevel)
	if err != nil {
		panic(fmt.Sprintf("init logger err %e", err))
	}
	cfg, err := internal.NewConfig(*projectPath, *buildPath, *sassDirPath, runMode)
	if err != nil {
		log.Info(err.Error())
		return
	}
	log.Info("watching project directory", zap.String("value", cfg.ProjectDir))
	time.Sleep(time.Second * 3)

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
	newMarkBus := pool.GetNewMarkOutBus()
	removeMarkBus := pool.GetNewMarkOutBus()
	go parser.Watch(ctx, cancel, compilerChan, newMarkBus, removeMarkBus)

	// start compilerChan worker pool

	for {
		select {
		case <-ctx.Done():
			log.Info("calceus watch was gracefully shutdown")
			close(compilerChan)
			close(done)
			return

		case <-done:
			log.Info("shutting down calceus watch...")
			cancel()
		}
	}
}

func initLogger(logLevel string) (*zap.Logger, error) {
	var levelValue zapcore.Level

	switch logLevel {
	case LogLevelInfo:
		levelValue = zapcore.InfoLevel

	case LogLevelDebug:
		levelValue = zapcore.DebugLevel

	default:
		return nil, errors.New("undefined log level")
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(levelValue)
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC822)
	return cfg.Build()
}
