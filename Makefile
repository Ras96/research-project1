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
	@go build ./...

lint:
	@go vet ./...
