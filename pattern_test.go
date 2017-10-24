package main

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		input string

		matched bool

		path string
		url  string
	}{
		{
			input:   "takonomura/git-get",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "https://github.com/takonomura/git-get.git",
		},
		{
			input:   "github.com/takonomura/git-get",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "https://github.com/takonomura/git-get.git",
		},
		{
			input:   "https://github.com/takonomura/git-get",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "https://github.com/takonomura/git-get.git",
		},
		{
			input:   "https://github.com/takonomura/git-get.git",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "https://github.com/takonomura/git-get.git",
		},
		{
			input:   "https://github.com/takonomura/git-get/",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "https://github.com/takonomura/git-get.git",
		},
		{
			input:   "https://github.com/takonomura/git-get.git/",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "https://github.com/takonomura/git-get.git",
		},
		{
			input:   "takonomura",
			matched: false,
		},
		{
			input:   "https://takonomura/git-get",
			matched: false,
		},
		{
			input:   "ssh://git@github.com:takonomura/git-get.git",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "ssh://git@github.com:takonomura/git-get.git",
		},
		{
			input:   "git@github.com:takonomura/git-get.git",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "ssh://git@github.com:takonomura/git-get.git",
		},
		{
			input:   "github.com:takonomura/git-get.git",
			matched: true,
			path:    "github.com/takonomura/git-get",
			url:     "ssh://github.com:takonomura/git-get.git",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			matched, repo := Match(test.input)
			if matched != test.matched {
				t.Fatalf("got matched = %v (want %v)", matched, test.matched)
			}
			if repo.Path != test.path {
				t.Fatalf("got path = %+v (want %+v)", repo.Path, test.path)
			}
			if repo.URL != test.url {
				t.Fatalf("got url = %+v (want %+v)", repo.URL, test.url)
			}
		})
	}
}
