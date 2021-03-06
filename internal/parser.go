package internal

import (
	"context"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Parser interface {
	AddCompiler(ctx context.Context, compilePath string)
	Watch(ctx context.Context, cancelFunc context.CancelFunc, compilerChan chan Compiler,
		newMarkChan chan<- string, removeMarkChan chan<- string)
}

type parser struct {
	log *zap.Logger

	// collection set for future processing files for each group by compiler
	m       sync.Mutex
	history History

	cfg Config
}

func (p *parser) AddCompiler(ctx context.Context, compilePath string) {
	//c := NewCompiler(ctx, p.log, compilePath, p.cfg)
	//p.watchingFiles[compilePath] = c
}

func (p *parser) GetDir() string {
	return p.cfg.ProjectDir
}

func (p *parser) GetBuildDir() string {
	return p.cfg.BuildDir
}

func (p *parser) Watch(ctx context.Context, cancelFunc context.CancelFunc, compilerChan chan Compiler,
	newMarkChan chan<- string, removeMarkChan chan<- string) {
	p.log.Info("start calceus parsing...")

	for {
		select {
		case <-time.After(WatchTimeout):
			p.history.Start()

			// walk through the directory tree
			err := p.walkByDir(p.GetDir())
			if err != nil {
				p.log.Error("get file names from root dir err", zap.Error(err))
			}

			p.logHistory()
			p.history.Commit()

			// get new history marks and start compilers
			newMarks, removeMarks := p.history.GetChanged()
			for _, value := range newMarks {
				p.log.Debug("send new mark", zap.Any("value", value))
				newMarkChan <- value
			}

			for _, value := range removeMarks {
				p.log.Debug("send remove mark", zap.Any("value", value))
				removeMarkChan <- value
			}
		}
	}
}

func (p *parser) addToHistory(filePath string) {
	p.m.Lock()
	p.history.Add(filePath)
	p.m.Unlock()
}

func (p *parser) logHistory() {
	p.history.LogInfo()
}

func (p *parser) walk(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		p.log.Error("read dir err", zap.Error(err))
		return err
	}

	for _, fileName := range files {
		if fileName.IsDir() && fileName.Name() == "node_modules" {
			continue
		} else if fileName.IsDir() {
			err := p.walk(path.Join(dir, fileName.Name()))
			if err != nil {
				p.log.Error("walk dir err", zap.Error(err))
			}
		} else {
			ok := p.isSASSPublicFile(fileName.Name())
			if !ok {
				p.addToHistory(path.Join(dir, fileName.Name()))
			}
		}
	}
	return nil
}

func (p *parser) walkByDir(dir string) error {
	sassDir := path.Join(dir, p.cfg.SassDir)
	return p.walk(sassDir)
}

func (p *parser) isSASSPublicFile(path string) bool {
	result := strings.Split(path, PrivateSASSFileDelimiter)
	return len(result) == 2
}

func NewParser(cfg Config, log *zap.Logger) Parser {
	return &parser{
		log:     log,
		history: NewHistory(log),
		cfg:     cfg,
		m:       sync.Mutex{},
	}
}
