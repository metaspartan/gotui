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

screenshots:
	go run _tools/generate.go

list:
	go run _tools/list_widgets.go

sexy:
	go fmt ./...
	$$(go env GOPATH)/bin/gocyclo -over 15 .
	$$(go env GOPATH)/bin/ineffassign ./...
