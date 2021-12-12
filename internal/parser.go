package internal

import (
	"go.uber.org/zap"
)

type Parser interface {
	Parse()
}

type parser struct {
	sassFiles []string
	log       *zap.Logger
}

func (p *parser) Parse() {
	p.log.Info("start calceus parsing...")
	// walk through the directory tree
	// run compilation
}

func NewParser(cfg Config, log *zap.Logger) Parser {
	return &parser{
		log:       log,
		sassFiles: make([]string, 0),
	}
}
