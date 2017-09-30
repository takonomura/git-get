package main

import (
	"path/filepath"
)

type RepoInfo struct {
	Path   string
	URL    string
	Branch string
}

func (r RepoInfo) CloneCmd(baseDir string) []string {
	dir := filepath.Join(baseDir, filepath.FromSlash(r.Path))

	cmd := []string{"git", "clone", "--recursive", r.URL, dir}

	if r.Branch != "" {
		cmd = append(cmd, "-b", r.Branch)
	}

	return cmd
}
