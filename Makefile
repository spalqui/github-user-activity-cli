.PHONY:
all: build test

build:
	go build -o github-activity .

test:
	go test ./...