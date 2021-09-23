GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install

FILE     := keys
KEY :=$(file < $(FILE))

BINARY_NAME=goshleep

build:
	$(GOBUILD) -v ./cmd/goshleep

install:
	$(GOINSTALL) github.com/Secret-Society-Blanket/goshleep/cmd/goshleep

run:
	$(GORUN) ./cmd/goshleep -t $(KEY) -p +
