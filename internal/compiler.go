package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	path2 "path"
	"path/filepath"

	"go.uber.org/zap"
)

type Compiler interface {
	// Build starts sass compiler execution
	Build(projectPath string) error

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

func (c *compiler) Build(projectPath string) error {
	c.log.Info("building compiler", zap.String("file", c.path))

	sassBinary := filepath.Join(projectPath, "node_modules/.bin/sass")
	cmdOpts := fmt.Sprintf("%s:%s", c.path, c.buildPath)
	cmd := exec.Command(sassBinary, cmdOpts, "--watch", "--no-source-map") // "--style=compressed"

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
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

	// build path with scss -> css file name
	cssFileName := filepath.Base(path[:len(path)-4]) + "css"
	buildPath := path2.Join(cfg.BuildDir, cssFileName)

	return &compiler{
		log:        log,
		ctx:        compilerCtx,
		cancelFunc: cancelFunc,
		path:       path,
		buildPath:  buildPath,
		mode:       cfg.Mode,
	}
}
