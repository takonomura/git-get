package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	branch string
	dryRun bool
)

func init() {
	flag.StringVarP(&branch, "branch", "b", "", "Branch to clone")
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run")
}

func getGitPath() string {
	p := os.Getenv("GITPATH")
	if p != "" {
		return p
	}
	p = os.Getenv("GOPATH")
	if p == "" {
		p = os.Getenv("HOME")
	}
	if p == "" {
		return p
	}
	return filepath.Join(p, "src")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-b branch] [--dry-run] repo\n", os.Args[0])
		os.Exit(1)
	}

	gitPath := getGitPath()
	if gitPath == "" {
		fmt.Fprintln(os.Stderr, "Please set $GITPATH")
		os.Exit(1)
	}

	matched, repo := Match(flag.Arg(0))
	if !matched {
		fmt.Fprintln(os.Stderr, "Cannot parse specified repository")
		os.Exit(1)
	}

	if branch != "" {
		repo.Branch = branch
	}

	cmd := repo.CloneCmd(gitPath)

	if dryRun {
		fmt.Println("$ " + strings.Join(cmd, " "))
		return
	}

	c := exec.Command(cmd[0], cmd[1:]...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clone: %s", err)
	}
}
