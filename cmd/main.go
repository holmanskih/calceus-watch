package main

import (
	"fmt"
	"github.com/holmanskih/calceus-watch/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	projectPath = "/Users/holmanskih/Desktop/calceus-sass/test_data/scss"
	buildPath   = "/Users/holmanskih/Desktop/calceus-sass/test_data/build/"
)

func main() {
	cfg := internal.Config{
		ProjectPath: projectPath,
		BuildPath:   buildPath,
		Mode:        internal.ModeProduction,
	}

	log, err := initLogger()
	if err != nil {
		panic(fmt.Sprintf("init logger err %e", err))
	}

	parser := internal.NewParser(cfg, log)
	parser.Parse()
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
