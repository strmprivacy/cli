.PHONY: build clean

SHELL := /bin/bash

build:
	goreleaser --snapshot --skip-publish --rm-dist

zsh-completion:
	/bin/zsh -c 'strm completion zsh > "$${fpath[1]}/_strm"'

# for a speedier build than with goreleaser
source_files := $(shell find . -name "*.go")
strm: ${source_files} Makefile
	go build -o $@
clean:
	rm -f strm
