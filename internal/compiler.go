package internal

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

type Compiler interface {
	// Build starts sass compiler execution
	Build(ctx context.Context) error

	// Kill stops sass compiler execution
	Kill()
}

type compiler struct {
	path      string
	buildPath string
	mode      Mode
}

func (c *compiler) Kill() {
	// todo: add kill logic(concerned with context stop)
}

func (c *compiler) Build(ctx context.Context) error {
	args := c.getBuildCmdArgs()

	cmd := exec.CommandContext(ctx, "npx sass", args...)
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

func NewCompiler(path string, cfg Config) Compiler {
	return &compiler{
		path:      path, // todo
		buildPath: "",   // todo
		mode:      cfg.Mode,
	}
}
