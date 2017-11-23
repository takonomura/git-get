package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	flag "github.com/spf13/pflag"
)

var (
	level  int
	branch string

	root  bool
	print bool
	list  bool

	outputString string
	output       *template.Template
)

func init() {
	flag.IntVarP(&level, "level", "L", 3, "Descend only level directories deep")
	flag.StringVarP(&branch, "branch", "b", "", "Branch to clone")

	flag.BoolVar(&root, "root", false, "Print $GITPATH")
	flag.BoolVarP(&list, "list", "l", false, "List repositories")
	flag.BoolVarP(&print, "print", "p", false, "Print")

	flag.StringVarP(&outputString, "output", "o", "", "Output template")

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

func parseOutput() error {
	if outputString == "" {
		switch {
		case list:
			outputString = `{{ . }}`
		case print:
			outputString = `{{ join (.CloneCmd root) " " }}`
		default:
			outputString = `$ {{ join (.CloneCmd root) " " }}`
		}
	}
	outputString += "\n"

	var err error
	output, err = template.New("output").Funcs(template.FuncMap{
		"join":         strings.Join,
		"filepath":     filepath.FromSlash,
		"filepathJoin": filepath.Join,
		"root": func() string {
			return gitPath
		},
		"abs": func(s string) string {
			return filepath.Join(gitPath, filepath.FromSlash(s))
		},
	}).Parse(outputString)
	return err
}

func printListAll() error {
	gitPath := filepath.Clean(gitPath) + string(filepath.Separator)
	return filepath.Walk(gitPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return err
		}

		rel := strings.TrimPrefix(path, gitPath)

		_, err = os.Stat(filepath.Join(path, ".git"))
		if err == nil {
			output.Execute(os.Stdout, rel)
			return filepath.SkipDir
		}

		if !os.IsNotExist(err) {
			return err
		}

		l := len(strings.Split(rel, string(filepath.Separator)))
		if l >= level {
			return filepath.SkipDir
		}
		return nil
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
	return execute(cmd)
}

func main() {
	if gitPath == "" {
		fmt.Fprintln(os.Stderr, "Please set $GITPATH")
		os.Exit(1)
	}

	if err := parseOutput(); err != nil {
		fmt.Fprintln(os.Stderr, "Invalid output template")
		os.Exit(1)
	}

	if root {
		fmt.Println(gitPath)
		return
	}

	if list {
		err := printListAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to find repositories: %s\n", err)
			os.Exit(1)
		}
		return
	}

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--print] [-b branch] repo\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s --root\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s --list\n", os.Args[0])
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

	output.Execute(os.Stdout, repo)
	if print {
		return
	}

	if err := clone(repo); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clone the repository: %s\n", err)
	}
}
