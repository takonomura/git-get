package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

type RepoInfo struct {
	Path   string
	URL    string
	Branch string
}

func (r RepoInfo) Clone(baseDir string) error {
	dir := filepath.Join(baseDir, filepath.FromSlash(r.Path))

	cmdName := "git"
	cmdArgs := []string{"clone", r.URL, dir}
	if r.Branch != "" {
		cmdArgs = append(cmdArgs, "-b", r.Branch)
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
