package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cwd, _ := os.Getwd()
	pathPrefix := filepath.Join(cwd, "scripts/sass-compiler/test_data")

	sassPath := filepath.Join(pathPrefix, "node_modules/.bin/sass")
	inPath := filepath.Join(pathPrefix, "in.scss")
	outPath := filepath.Join(pathPrefix, "out.css")
	comileOpts := fmt.Sprintf("%s:%s", inPath, outPath)
	cmd := exec.Command(sassPath, comileOpts, "--watch")

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
