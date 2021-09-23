GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install

BINARY_NAME=goshleep

build:
	$(GOBUILD) -v ./cmd/goshleep

install:
	$(GOINSTALL) github.com/Secret-Society-Blanket/goshleep/cmd/goshleep
