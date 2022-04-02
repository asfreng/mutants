##
## Challenge
##

.DEFAULT_GOAL := help

## build: Build challenge and put the binary in the bin directory
.PHONY: build
build:
	go build -o bin/challenge cmd/challenge/main.go

## run: Runs challenge app
.PHONY: run
run:
	go run cmd/challenge/main.go

## test: Run the tests
.PHONY: test
test:
	go test ./... -cover

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
