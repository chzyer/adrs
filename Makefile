export GOBIN := $(shell pwd)/bin
export PREFIX := github.com/chzyer/adrs

all: deps
	go install ./...

deps:
	godep restore

test:
	go test ./...

cover:
	go test ./... -cover
