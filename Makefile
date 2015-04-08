export GOBIN := $(shell pwd)/bin
export PREFIX := github.com/chzyer/adrs

all:
	go install ./...

test:
	go install ${PREFIX}/main/adrs-test
	go install ${PREFIX}/main/adrs-udp
