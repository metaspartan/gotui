.PHONY: run-examples test build

run-examples:
	@for file in _examples/*.go; do \
	  go run $$file; \
	  done;

test:
	go test ./_test/block_test.go
	go test ./_test/style_test.go

build:
	go build ./...
