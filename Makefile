.PHONY: build test clean all
.DEFAULT_GOAL := all

SHELL := /bin/bash

build:
	goreleaser --snapshot --skip-publish --rm-dist

test:
	go test ./test

zsh-completion:
	/bin/zsh -c 'strm completion zsh > "$${fpath[1]}/_strm"'

# for a speedier build than with goreleaser
source_files := $(shell find . -name "*.go")

targetVar := streammachine.io/strm/cmd.CommandName

target := dstrm

ldflags := -X '${targetVar}=${target}'

${target}: ${source_files} Makefile
	go build -ldflags="${ldflags}" -o $@

clean:
	rm -f ${target}

all: ${target}
