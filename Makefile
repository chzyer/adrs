export GOBIN := $(shell pwd)/build
export PREFIX := github.com/chzyer/adrs

all:
	go install ./...

test:
	go install ${PREFIX}/bin/adrs-test
