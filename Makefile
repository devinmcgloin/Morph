GO_FILES := $(shell find . -iname '*.go' -type f | grep -v /vendor/) # All the .go files, excluding vendor/

install:
	go install ./...

test:
	test -z $(shell gofmt -s -l $(GO_FILES))         # Fail if a .go file hasn't been formatted with gofmt
	go test -v -race ./...                   # Run all the tests with the race detector enabled
	go vet ./...                             # go vet is the official Go static analyzer
	megacheck ./...

vet:
	go vet ./...
