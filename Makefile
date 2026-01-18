GO      	:= go
GOBUILD 	:= $(GO) build
GOTEST  	:= $(GO) test
GOCLEAN		:= $(GO) clean

GEN_FOLDER	:= gen
PROTO_FOLDER = proto/music

BINARY  	:= bin/music-service

.PHONY: all build test clean run
	
all: run

build: gen
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
	protoc --proto_path=${PROTO_FOLDER} --go_out=. --go-grpc_out=. ${PROTO_FOLDER}/models.proto ${PROTO_FOLDER}/service.proto 
