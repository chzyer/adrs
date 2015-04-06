export GOBIN := $(shell pwd)/build
all:
	go get -u gopkg.in/logex.v1
	go install ./...
