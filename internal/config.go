package internal

type Mode uint

const (
	ModeProduction Mode = iota
	ModeDevelopment
)

type Config struct {
	ProjectPath string
	BuildPath   string
	Mode        Mode
}
