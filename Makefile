.PHONY: build test clean all
.DEFAULT_GOAL := all

SHELL := /bin/bash

build:
	goreleaser --snapshot --skip-publish --rm-dist

zsh-completion:
	/bin/zsh -c 'strm completion zsh > "$${fpath[1]}/_strm"'

# for a speedier build than with goreleaser
source_files := $(shell find . -name "*.go")

targetVar := streammachine.io/strm/pkg/common.RootCommandName

target := dstrm

ldflags := -X '${targetVar}=${target}' -X streammachine.io/strm/pkg/cmd.Version=local -X streammachine.io/strm/pkg/cmd.GitSha=local -X streammachine.io/strm/pkg/cmd.BuiltOn=local

dist/${target}: ${source_files} Makefile
	go build -ldflags="${ldflags}" -o $@ ./cmd/strm

clean:
	rm -f ${target}

test: dist/${target}
	go clean -testcache
	go test ./test -v

all: dist/${target}
