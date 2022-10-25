SHELL:=/bin/bash
CMDS := $(shell grep -E -o "^[a-z_-]+" ./Makefile)

.PHONY: $(CMDS)

commands:
	@echo "Available commands:"
	@echo "-------------------"
	@echo $(CMDS) | sed 's/\s/\n/g'

all: clean mod build
all-debug: clean mod build-debugb

clean:
	@go clean

mod:
	@go mod tidy

build:
	@go build .

build-debug:
	@go build -tags debug .

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix ./...

download-corpus:
	@mkdir -p corpus
	@wget -O corpus/corpus.zip "https://sites.google.com/site/dialoguebreakdowndetection/chat-dialogue-corpus/projectnextnlp-chat-dialogue-corpus.zip"
	@unzip corpus/corpus.zip -d corpus
	@rm corpus/corpus.zip

marp:
	@docker run --rm -v $(shell pwd):/home/marp/app/ -e LANG=$LANG -e MARP_USER="$(shell id -u):$(shell id -g)" marpteam/marp-cli -o ./dist/index.html ./docs/slide.md
