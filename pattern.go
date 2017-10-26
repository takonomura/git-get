package main

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
)

type Pattern func(s string) (matched bool, repo RepoInfo)

func RegexpPattern(pattern string, parser func([]string) RepoInfo) Pattern {
	re := regexp.MustCompile(pattern)
	return func(s string) (matched bool, repo RepoInfo) {
		matches := re.FindAllStringSubmatch(s, 1)
		if len(matches) != 1 {
			return
		}
		matched = true
		parts := matches[0]
		repo = parser(parts)
		return
	}
}

var Patterns = []Pattern{
	RegexpPattern(`^(?:(?:https://)?([a-zA-Z0-9-.]+)/)?([a-zA-Z0-9-_.]+)/([a-zA-Z0-9-_.]+?)(?:\.git)?/?$`, func(parts []string) RepoInfo {
		if parts[1] == "" {
			parts[1] = "github.com"
		}
		return RepoInfo{
			Path: filepath.Join(parts[1:]...),
			URL:  "https://" + path.Join(parts[1:]...) + ".git",
		}
	}),
	RegexpPattern(`^(?:ssh://)?(?:([a-z0-9_]+)@)?([a-zA-Z0-9-.]+):([a-zA-Z0-9-_./]+?).git$`, func(parts []string) RepoInfo {
		path := filepath.FromSlash(parts[3])
		path = filepath.Join(parts[2], path)
		var url string
		if parts[1] == "" {
			url = fmt.Sprintf("ssh://%s:%s.git", parts[2], parts[3])
		} else {
			url = fmt.Sprintf("ssh://%s@%s:%s.git", parts[1], parts[2], parts[3])
		}
		return RepoInfo{
			Path: path,
			URL:  url,
		}
	}),
}

func Match(s string) (matched bool, repo RepoInfo) {
	for _, pattern := range Patterns {
		if matched, repo := pattern(s); matched {
			return matched, repo
		}
	}
	return
}
