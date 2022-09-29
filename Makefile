.PHONY: build test clean all
.DEFAULT_GOAL := all

SHELL := /bin/bash

build:
	goreleaser --snapshot --skip-publish --rm-dist

zsh-completion:
	/bin/zsh -c 'strm completion zsh > "$${fpath[1]}/_strm"'

# for a speedier build than with goreleaser
source_files := $(shell find . -name "*.go")

targetVar := strmprivacy/strm/pkg/common.RootCommandName

target := dstrm

ldflags := -X '${targetVar}=${target}' -X strmprivacy/strm/pkg/cmd.Version=local -X strmprivacy/strm/pkg/common.GitSha=local -X strmprivacy/strm/pkg/common.BuiltOn=local

dist/${target}: ${source_files} Makefile
	go build -ldflags="${ldflags}" -o $@ ./cmd/strm

clean:
	rm -f dist/${target}

test: dist/${target}
	go clean -testcache
	go test ./test -v

all: dist/${target}
