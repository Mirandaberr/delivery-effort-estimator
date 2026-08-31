.PHONY: fmt vet test build

fmt:
	gofmt -l .

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./...
