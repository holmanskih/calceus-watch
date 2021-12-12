package internal

import (
	"context"
	"fmt"
	"os/exec"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

type Compiler interface {
	// Build starts sass compiler execution
	Build() error

	// Kill stops sass compiler execution
	Kill()
}

type compiler struct {
	log *zap.Logger

	ctx        context.Context
	cancelFunc context.CancelFunc
	path       string
	buildPath  string
	mode       Mode
}

func (c *compiler) Kill() {
	c.log.Info("killing compiler", zap.String("file", c.path))
	c.cancelFunc()
}

func (c *compiler) Build() error {
	c.log.Info("building compiler", zap.String("file", c.path))

	args := c.getBuildCmdArgs()

	cmd := exec.CommandContext(c.ctx, "npx sass", args...)
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "build ")
	}

	return nil
}

func (b *compiler) getBuildCmdArgs() []string {
	fileInfo := fmt.Sprintf("%s:%s", b.path, b.buildPath)

	if b.mode == ModeProduction {
		return []string{"--no-source-map", fileInfo, "--watch", "--style=compressed"}
	}

	return []string{"--no-source-map", fileInfo, "--watch"}
}

func NewCompiler(ctx context.Context, log *zap.Logger, path string, cfg Config) Compiler {
	compilerCtx, cancelFunc := context.WithCancel(ctx)

	return &compiler{
		log:        log,
		ctx:        compilerCtx,
		cancelFunc: cancelFunc,
		path:       path, // todo
		buildPath:  "",   // todo
		mode:       cfg.Mode,
	}
}
