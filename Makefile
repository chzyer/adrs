export GOBIN := $(shell pwd)/build
all:
	go get gopkg.in/logex.v1
	go install ./...
