export GOBIN := $(shell pwd)/bin
export PREFIX := github.com/chzyer/adrs

all:
	go install ./...

test:
	go test ./...

cover:
	go test ./... -cover

cover-func:
	./goverall.sh func
