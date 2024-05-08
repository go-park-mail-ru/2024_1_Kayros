BINARY_NAME = main
ENTRYPOINT = cmd/main/main.go

all: build test

build:
	go build -o ${BINARY_NAME} ${ENTRYPOINT}

test:
	go test -v main.go

run: build
	${BINARY_NAME}

dep:
	go fmt ./...
	go mod tidy -v

clean:
	go clean
	rm ${BINARY_NAME}