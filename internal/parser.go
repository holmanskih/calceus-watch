package internal

import (
	"io/fs"
	"path/filepath"

	"go.uber.org/zap"
)

type Parser interface {
	Parse()
}

type parser struct {
	log       *zap.Logger
	sassFiles []string

	cfg Config
}

func (p *parser) GetDir() string {
	return p.cfg.Dir
}

func (p *parser) GetBuildDir() string {
	return p.cfg.BuildDir
}

func (p *parser) Parse() {
	p.log.Info("start calceus parsing...")

	// walk through the directory tree
	err := p.getFileNamesFromDir(p.GetDir())
	if err != nil {
		p.log.Error("get file names from root dir err", zap.Error(err))
	}

	// todo: run compilation
}

func (p *parser) getFileNamesFromDir(dir string) error {
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			p.log.Error("dir walk err", zap.Error(err))
			return err
		}
		p.log.Debug("found file", zap.String("fileName", info.Name()))
		return nil
	})
	return nil
}

func NewParser(cfg Config, log *zap.Logger) Parser {
	return &parser{
		log:       log,
		sassFiles: make([]string, 0),
		cfg:       cfg,
	}
}
