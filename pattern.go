package main

import (
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
	RegexpPattern(`^(?:(?:https://)?([a-zA-Z0-9-.]+)/)?([a-zA-Z0-9-_.]+)/([a-zA-Z0-9-_.]+?)(?:\.git)?$`, func(parts []string) RepoInfo {
		if parts[1] == "" {
			parts[1] = "github.com"
		}
		return RepoInfo{
			Path: filepath.Join(parts[1:]...),
			URL:  "https://" + path.Join(parts[1:]...) + ".git",
		}
	}),
}

func Match(s string) (matched bool, repo RepoInfo) {
	for _, pattern := range Patterns {
		if matched, repo := pattern(s); matched {
			return matched, repo
		}
	}
	return false, repo
}
