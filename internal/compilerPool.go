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

func (p *compilerPool) GetNewMarkOutBus() chan<- string {
	return p.newMarkChan
}

func (p *compilerPool) GetRemoveMarkOutBus() chan<- string {
	return p.removeMarkChan
}

func (p *compilerPool) GetRemoveCompilerOutBus() chan<- string {
	return p.removeMarkChan
}

func (p *compilerPool) Run(ctx context.Context, cfg Config) {
	for {
		select {
		case <-ctx.Done():
			p.log.Debug("compiler pool ctx is done")
			//close(p.removeMarkChan)
			close(p.newMarkChan)

		case mark, ok := <-p.newMarkChan:
			if !ok {
				p.log.Debug("read from new mark chan err")
			}

			p.log.Debug("receive new mark of type [new]", zap.Any("value", mark))
			go func(ctx context.Context) {
				comp := NewCompiler(ctx, p.log, mark, cfg)
				comp.Build(cfg.ProjectDir)
			}(ctx)

			//case mark, ok := <-p.removeMarkChan:
			//	if !ok {
			//		p.log.Debug("read from remove mark chan err")
			//	}
			//	p.log.Info("receive new mark of type [remove]", zap.Any("value", mark))
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
