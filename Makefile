GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=aprsweb-native
BINARY_ARM=aprsweb
VERSION := $(shell git rev-parse --short HEAD)
LDFLAGS=-w -s -X 'main.Version=$(VERSION)'
all: test build arm 
build:
	CGO_ENABLED=0 $(GOBUILD)  -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) -v ./cmd/aprsweb/aprsweb.go
arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD)  -ldflags="$(LDFLAGS)" -o $(BINARY_ARM) -v ./cmd/aprsweb/aprsweb.go	
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_ARM)

