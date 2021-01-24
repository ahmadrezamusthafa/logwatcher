build:
	@echo "Building binary file" && \
	go build -o logwatcher main.go && \
	echo "Done"

SYMLINK = /mnt/c/Users/Reza/go
	export GOPATH = $(shell go env GOPATH):$(SYMLINK):$(shell pwd)/vendor:$(shell pwd)
	export GOBIN = $(SYMLINK)/bin
	export PATH = $(shell printenv PATH):$(GOPATH):$(GOBIN)

build_windows:
	@echo "Building binary file" && \
	env GOOS=windows GOARCH=386 go build -o logwatcher.exe main.go && \
	echo "Done"