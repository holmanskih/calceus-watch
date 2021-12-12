package internal

import (
	"io/fs"
	"path/filepath"
	"strings"

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

	// todo: run compilation(temporary solution)

}

func (p *parser) walkByDir(path string, info fs.FileInfo, err error) error {
	if err != nil {
		p.log.Error("dir walk err", zap.Error(err))
		return err
	}

	if info.IsDir() {
		return nil
	}
	ok := p.isSASSPublicFile(info.Name())
	if !ok {
		//p.sassFiles = append(p.sassFiles, info.Name())
		p.log.Debug("found file", zap.String("fileName", info.Name()))
	}
	return nil
}

func (p *parser) getFileNamesFromDir(dir string) error {
	err := filepath.Walk(dir, p.walkByDir)
	if err != nil {
		p.log.Error("dir walk err", zap.Error(err))
	}

	return nil
}

func (p *parser) isSASSPublicFile(path string) bool {
	result := strings.Split(path, PrivateSASSFileDelimiter)
	return len(result) == 2
}

func NewParser(cfg Config, log *zap.Logger) Parser {
	return &parser{
		log:       log,
		sassFiles: make([]string, 0),
		cfg:       cfg,
	}
}
