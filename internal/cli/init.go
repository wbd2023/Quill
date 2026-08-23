package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

/* ---------------------------------------- Init Command ---------------------------------------- */

type initCmd struct {
	Root   string `kong:"name=repository-root,help=target directory (current directory when omitted)"`
	Preset string `kong:"default=minimal,enum='minimal',help=preset: minimal"`
}

func (c *initCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	root, err := resolveInitTarget(c.Root)
	if err != nil {
		return runner.reportCommandError("init", report.FormatText, err)
	}

	result, err := engine.Init(ctx, root, c.Preset)
	if err != nil {
		return runner.reportCommandError("init", report.FormatText, err)
	}
	if err := report.WriteInit(runner.stdout, result); err != nil {
		return runner.reportCommandError("init", report.FormatText, err)
	}

	return 0
}

// resolveInitTarget resolves the directory init writes into. Unlike repository-root discovery,
// init targets a directory that may not yet be a repository: an explicit path wins, otherwise the
// current working directory is used.
func resolveInitTarget(path string) (target string, err error) {
	if path != "" {
		return filepath.Abs(path)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	return workingDirectory, nil
}
