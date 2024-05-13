BINARY_NAME = main
ENTRYPOINT = cmd/main/main.go

all: build test

build:
	go build -o ${BINARY_NAME} ${ENTRYPOINT}

test:
	go test -v main.go

tidy:
	go fmt ./...
	go mod tidy -v

clean:
	docker rm gateway
	docker image rm dev-compose-gateway

generate:
	go generate 

pb:
	protoc -I proto proto/**/*.proto --go_out=gen/go --go-grpc_out=gen/go