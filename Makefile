GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
BIN_NAME=rsoi


all: build

build: fmt
	rm -f $(BIN_NAME)
	$(GOBUILD) -o $(BIN_NAME)

test:
	GOCACHE=off $(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BIN_NAME)

run: build
	./$(BIN_NAME)

fmt:
	$(GOFMT) ./...
