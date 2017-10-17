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

	flag.Parse()
}

var gitPath = getGitPath()

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

func printListAll() error {
	gitPath := filepath.Clean(gitPath) + string(filepath.Separator)
	return filepath.Walk(gitPath, func(path string, info os.FileInfo, err error) error {
		rel := strings.TrimPrefix(path, gitPath)
		if _, err := os.Stat(filepath.Join(path, ".git")); err != nil {
			if len(strings.Split(rel, string(filepath.Separator))) >= 3 {
				return filepath.SkipDir
			}
			return nil
		}
		if fullPath {
			fmt.Println(path)
		} else {
			fmt.Println(rel)
		}
		return filepath.SkipDir
	})
}

func execute(cmd []string) error {
	binary, err := exec.LookPath(cmd[0])
	if err != nil {
		return err
	}
	return syscall.Exec(binary, cmd, os.Environ())
}

func clone(repo RepoInfo) error {
	cmd := repo.CloneCmd(gitPath)

	fmt.Println("$ " + strings.Join(cmd, " "))

	if dryRun {
		return nil
	}

	return execute(cmd)
}

func main() {
	if gitPath == "" {
		fmt.Fprintln(os.Stderr, "Please set $GITPATH")
		os.Exit(1)
	}

	if list {
		err := printListAll()
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

	if err := clone(repo); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clone the repository: %s\n", err)
	}
}
