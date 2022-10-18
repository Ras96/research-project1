SHELL:=/bin/bash
CMDS := $(shell grep -E -o "^[a-z_-]+" ./Makefile)

.PHONY: $(CMDS)

commands:
	@echo "Available commands:"
	@echo "-------------------"
	@echo $(CMDS) | sed 's/\s/\n/g'

all: clean mod build

clean:
	@go clean

mod:
	@go mod tidy

build:
	@go build .

lint:
	@go vet ./...

download-corpus:
	@mkdir -p corpus
	@wget -O corpus/corpus.zip "https://sites.google.com/site/dialoguebreakdowndetection/chat-dialogue-corpus/projectnextnlp-chat-dialogue-corpus.zip"
	@unzip corpus/corpus.zip -d corpus
	@rm corpus/corpus.zip
