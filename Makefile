BINARY=gaelv
PACKAGE=gaelv

all: test build

test:
	go test -v ./...

build:
	go-bindata -pkg $(PACKAGE) static/ templates/
	go build -o $(BINARY) cmd/gaelv/main.go

clean:
	rm $(BINARY)
