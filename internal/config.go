package internal

type Mode uint

const (
	ModeProduction Mode = iota
	ModeDevelopment
)

type Config struct {
	Dir      string
	BuildDir string
	Mode     Mode
}

const PrivateSASSFileDelimiter = "_"
