.PHONY: build

SHELL := /bin/bash

build:
	goreleaser --snapshot --skip-publish --rm-dist

zsh-completion:
	/bin/zsh -c 'strm completion zsh > "$${fpath[1]}/_strm"'

