build:
	go build -o . ./...

.PHONY: test
test:
	go test ./...