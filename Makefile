.PHONY: test test-cover gen gen-ic fmt

test:
	go test -v -cover ./... --count=5

fmt:
	go mod tidy
	gofmt -s -w .
	goarrange run -r .
	golangci-lint run ./...
