GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOTOOL=$(GOCMD) tool
BIN_DIR=release
BIN_NAME=rsoi


all: build

build: fmt
	rm -rf $(BIN_DIR)
	mkdir $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BIN_NAME)

test:
	GOCACHE=off $(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)

run: build
	./$(BIN_DIR)/$(BIN_NAME)

fmt:
	$(GOFMT) ./...

cover:
	$(GOTEST) -coverprofile cp.out ./...
	$(GOTOOL) cover -html=cp.out
