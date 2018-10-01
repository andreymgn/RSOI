GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOTOOL=$(GOCMD) tool


all: build

test:
	GOCACHE=off $(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)

fmt:
	$(GOFMT) ./...

cover:
	$(GOTEST) -coverprofile cp.out ./...
	$(GOTOOL) cover -html=cp.out

proto:
	for f in services/**/proto/*.proto; do \
		protoc --go_out=plugins=grpc:. $$f; \
		echo compiled: $$f; \
	done

build:
	$(GOBUILD) ./cmd/...

debug:
	./RSOI -service post -port 8081 -conn "user=postgres password=secret dbname=postgres sslmode=disable port=5433" -jaeger-addr localhost:6831 &
	./RSOI -service comment -port 8082 -conn "user=postgres password=secret dbname=postgres sslmode=disable port=5434" -jaeger-addr localhost:6831 &
	./RSOI -service poststats -port 8083 -conn "user=postgres password=secret dbname=postgres sslmode=disable port=5435" -jaeger-addr localhost:6831 &
	./RSOI -service api -port 8080 -post-server "localhost:8081" -comment-server "localhost:8082" -post-stats-server "localhost:8083" -jaeger-addr localhost:6831 &
