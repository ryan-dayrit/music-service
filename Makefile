GO      	:= go
GOBUILD 	:= $(GO) build
GOTEST  	:= $(GO) test
GOCLEAN		:= $(GO) clean
BINARY  	:= postgres-crud
GEN_FOLDER	:= gen

.PHONY: all build test clean run
	
all: clean gen build

build:
	$(GOBUILD) -o $(BINARY) .

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY)
	rm -rf ${GEN_FOLDER}

run: build
	./$(BINARY)

gen: clean
	mkdir -p ${GEN_FOLDER}
	protoc --proto_path=./proto/album --go_out=$(GEN_FOLDER) --go-grpc_out $(GEN_FOLDER) ./proto/album/models.proto ./proto/album/service.proto 

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
