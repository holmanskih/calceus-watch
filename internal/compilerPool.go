package internal

import (
	"context"

	"go.uber.org/zap"
)

type compilerPool struct {
	log            *zap.Logger
	newMarkChan    chan string
	removeMarkChan chan string
}

func (p *compilerPool) GetNewCompilerOutBus() chan<- string {
	return p.newMarkChan
}

func (p *compilerPool) GetRemoveCompilerOutBus() chan<- string {
	return p.removeMarkChan
}

func (p *compilerPool) Run(ctx context.Context, cfg Config) {
	for {
		select {
		case <-ctx.Done():
			p.log.Debug("compiler pool ctx is done")

		case mark, ok := <-p.newMarkChan:
			if !ok {
				p.log.Debug("read from new mark chan err")
			}

			p.log.Info("receive new mark", zap.Any("value", mark))
			go func(ctx context.Context) {
				comp := NewCompiler(ctx, p.log, mark, cfg)
				comp.Build(cfg.ProjectDir)
			}(ctx)
		}
	}
}

func NewCompilerPool(log *zap.Logger) *compilerPool {
	return &compilerPool{
		log:            log,
		newMarkChan:    make(chan string),
		removeMarkChan: make(chan string),
	}
}
