package internal

import "time"

type Mode uint

const (
	ModeProduction Mode = iota
	ModeDevelopment

	WatchTimeout = time.Second * 2
)

type Config struct {
	ProjectDir string
	BuildDir   string
	Mode       Mode
}

const PrivateSASSFileDelimiter = "_"
