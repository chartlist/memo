# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

all: clean build

build:
	$(GOBUILD) -o ./bin/memo -v cmd/main/main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -f ./bin/memo

linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD)  -o ./bin/memo-linux -v cmd/main/main.go
