build:
	go build -o . ./...

.PHONY: test
test:
	go test ./...

.PHONY: dist
dist:
	node ./scripts/release $(TAG)