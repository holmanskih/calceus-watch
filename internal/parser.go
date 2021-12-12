package internal

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

const watchTimeout = time.Second * 1

type Parser interface {
	Watch(ctx context.Context, cancelFunc context.CancelFunc, compilerChan chan Compiler)
}

type parser struct {
	log *zap.Logger

	m               sync.Mutex
	watchingFiles   []string
	pretendingFiles []string
	cfg             Config
}

func (p *parser) GetDir() string {
	return p.cfg.Dir
}

func (p *parser) GetBuildDir() string {
	return p.cfg.BuildDir
}

func (p *parser) Watch(ctx context.Context, cancelFunc context.CancelFunc, compilerChan chan Compiler) {
	p.log.Info("start calceus parsing...")

	// todo: ass smart system for file transformation watch
	for {
		select {
		case <-time.After(watchTimeout):
			// walk through the directory tree
			err := p.walkByDir(p.GetDir())
			if err != nil {
				p.log.Error("get file names from root dir err", zap.Error(err))
			}

			p.m.Lock()
			ok := p.isWatchingFilesChanged(p.watchingFiles, p.pretendingFiles)
			if !ok {
				p.watchingFiles = p.pretendingFiles
			}
			p.m.Unlock()

			p.getWatchFilesInfo()
		}
	}
	// todo: run compilation(temporary solution)
}

func (p *parser) getWatchFilesInfo() {
	temp := make([]string, 0)
	for _, filePath := range p.watchingFiles {
		temp = append(temp, filepath.Base(filePath))
	}
	p.log.Debug("watching sass file", zap.Any("file", temp))
}

func (p *parser) walk(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		p.log.Error("read dir err", zap.Error(err))
		return err
	}

	for _, fileName := range files {
		if fileName.IsDir() {
			err := p.walk(path.Join(dir, fileName.Name()))
			if err != nil {
				p.log.Error("walk dir err", zap.Error(err))
			}
		} else {
			ok := p.isSASSPublicFile(fileName.Name())
			if !ok {
				p.addPretendingFile(dir, fileName.Name())
			}
		}
	}
	return nil
}

func (p *parser) addPretendingFile(dir, file string) {
	p.m.Lock()
	p.pretendingFiles = append(p.pretendingFiles, path.Join(dir, file))
	p.m.Unlock()
}

func (p *parser) walkByDir(dir string) error {
	p.pretendingFiles = make([]string, 0)
	return p.walk(dir)
}

func (p *parser) isWatchingFilesChanged(old, new []string) bool {
	if len(old) != len(new) {
		return false
	}
	for i, value := range old {
		if value != new[i] {
			return false
		}
	}
	return true
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
		m:             sync.Mutex{},
	}
}
