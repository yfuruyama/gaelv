NAME=gaelv

all: test build

test:
	go test -v ./...

build:
	go-bindata -pkg $(NAME) static/ templates/
	go build -o $(NAME) cmd/gaelv/main.go

install:
	go install github.com/addsict/$(NAME)/cmd/$(NAME)

clean:
	rm $(NAME)
