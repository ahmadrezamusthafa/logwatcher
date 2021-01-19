SYMLINK = /mnt/c/Users/Reza/go
export GOPATH = $(shell go env GOPATH):$(SYMLINK):$(shell pwd)/vendor:$(shell pwd)
export GOBIN = $(SYMLINK)/bin
export PATH = $(shell printenv PATH):$(GOPATH):$(GOBIN)

build:
	@echo "Building binary file" && \
	env GOOS=windows GOARCH=386 go build -o ./bin/logwatcher.exe main.go && \
	echo "Done"