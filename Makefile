.PHONY: test fmt

test:
	go test -v -cover ./...

fmt:
	go mod tidy
	gofmt -s -w .
	goarrange run -r .
	golangci-lint run ./...

report:
	TEST_SUITE_REPORT=true go test ./...
