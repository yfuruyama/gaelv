BINARY=gaelv

all: test build

test:
	go test -v ./...

build:
	go build -o $(BINARY) cmd/gaelv/main.go

clean:
	rm $(BINARY)
