export GOBIN := $(shell pwd)/build

all:
	go install ./...
