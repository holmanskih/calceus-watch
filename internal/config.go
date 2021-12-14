package internal

import (
	"time"

	"github.com/pkg/errors"
)

type Mode uint

const (
	ModeProduction Mode = iota
	ModeDevelopment

	WatchTimeout = time.Second * 2
)

type Config struct {
	ProjectDir string
	BuildDir   string
	SassDir    string
	Mode       Mode
}

func NewConfig(projectDir, buildDir, sassDir string, mode Mode) (Config, error) {
	if projectDir == "" {
		return Config{}, errors.New("project directory path was not specified")
	}

	if buildDir == "" {
		return Config{}, errors.New("build directory path was not specified")
	}

	if sassDir == "" {
		return Config{}, errors.New("sass directory path was not specified")
	}

	return Config{
		ProjectDir: projectDir,
		BuildDir:   buildDir,
		SassDir:    sassDir,
		Mode:       mode,
	}, nil
}

const PrivateSASSFileDelimiter = "_"
