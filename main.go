package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	flag "github.com/spf13/pflag"
)

var (
	branch string
	dryRun bool

	list     bool
	fullPath bool
)

func init() {
	flag.StringVarP(&branch, "branch", "b", "", "Branch to clone")
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run")

	flag.BoolVarP(&list, "list", "l", false, "List repositories")
	flag.BoolVarP(&fullPath, "path", "p", false, "Print full path")
}

func absPath(p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Join(os.Getenv("HOME"), p)
	}
	return p
}

func getGitPath() string {
	if p := os.Getenv("GITPATH"); p != "" {
		return absPath(p)
	}
	if p := os.Getenv("GOPATH"); p != "" {
		return absPath(filepath.Join(p, "src"))
	}
	if p := os.Getenv("HOME"); p != "" {
		return filepath.Join(p, "src")
	}
	return ""
}

func printList(gitPath string, fullPath bool) error {
	gitPath = filepath.Clean(gitPath) + string(filepath.Separator)
	return filepath.Walk(gitPath, func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(filepath.Join(path, ".git")); err != nil {
			if len(strings.Split(path, string(filepath.Separator))) >= 3 {
				return filepath.SkipDir
			}
			return nil
		}
		rel := strings.TrimPrefix(path, gitPath)
		if fullPath {
			fmt.Println(path)
		} else {
			fmt.Println(rel)
		}
		return filepath.SkipDir
	})
}

func main() {
	flag.Parse()

	gitPath := getGitPath()
	if gitPath == "" {
		fmt.Fprintln(os.Stderr, "Please set $GITPATH")
		os.Exit(1)
	}

	if list {
		err := printList(gitPath, fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to find repositories: %s", err)
			os.Exit(1)
		}
		return
	}

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-b branch] [--dry-run] repo\n", os.Args[0])
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

	if err := execute(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to exec command: %s\n", err)
	}
}

func execute(cmd []string) error {
	binary, err := exec.LookPath(cmd[0])
	if err != nil {
		return err
	}
	return syscall.Exec(binary, cmd, os.Environ())
}
