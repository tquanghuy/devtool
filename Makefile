.PHONY: build run test clean

build:
	go build -o devtool ./cmd/devtool

run: build
	./devtool

test:
	go test ./...

clean:
	rm -f devtool
