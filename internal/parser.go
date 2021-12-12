package internal

import (
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
)

type Parser interface {
	Parse()
}

type parser struct {
	log           *zap.Logger
	watchingFiles []string

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
	err := p.walkByDir(p.GetDir())
	if err != nil {
		p.log.Error("get file names from root dir err", zap.Error(err))
	}

	for _, filePath := range p.watchingFiles {
		p.log.Debug("watching sass file", zap.Any("file", filePath))
	}
	// todo: run compilation(temporary solution)
}

func (p *parser) walkByDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		p.log.Error("read dir err", zap.Error(err))
	}

	for _, fileName := range files {
		if fileName.IsDir() {
			err := p.walkByDir(path.Join(dir, fileName.Name()))
			if err != nil {
				p.log.Error("walk dir err", zap.Error(err))
			}
		} else {
			ok := p.isSASSPublicFile(fileName.Name())
			if !ok {
				p.watchingFiles = append(p.watchingFiles, path.Join(dir, fileName.Name()))
			}
		}
	}
	return nil
}

func (p *parser) isSASSPublicFile(path string) bool {
	result := strings.Split(path, PrivateSASSFileDelimiter)
	return len(result) == 2
}

func NewParser(cfg Config, log *zap.Logger) Parser {
	return &parser{
		log:           log,
		watchingFiles: make([]string, 0),
		cfg:           cfg,
	}
}
