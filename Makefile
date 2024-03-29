.PHONY: build test clean all
.DEFAULT_GOAL := all

SHELL := /bin/bash

build:
	goreleaser --snapshot --skip-publish --clean

zsh-completion:
	strm completion zsh > "$${fpath[1]}/_strm"

# for a speedier build than with goreleaser
source_files := $(shell find . -name "*.go")

targetVar := strmprivacy/strm/pkg/common.RootCommandName

target := dstrm

ldflags := -X '${targetVar}=${target}' -X strmprivacy/strm/pkg/cmd.Version=local -X strmprivacy/strm/pkg/common.GitSha=local -X strmprivacy/strm/pkg/common.BuiltOn=local

dist/${target}: ${source_files} Makefile
	go build -ldflags="${ldflags}" -o $@ ./cmd/strm

clean:
	rm -f dist/${target}

# Make sure the .env containing all `STRM_TEST_*` variables is present in the ./test directory
# godotenv loads the .env file from that directory when running the tests
test: dist/${target}
	go clean -testcache
	go test ./test -v

all: dist/${target}
