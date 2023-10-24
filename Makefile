.PHONY: test fmt report download

# Runs the tests.
test:
	go test -v -cover ./...

# Formats the code and runs the linter.
fmt:
	go mod tidy
	gofmt -s -w .
	goarrange run -r .
	golangci-lint run ./...

# Generates a test report for the W3C test suite.
report:
	TEST_SUITE_REPORT=true go test ./...

# Downloads the latest test suite from the W3C website.
download:
	rm ntriples/testdata/suite/* && curl -s -L https://www.w3.org/2013/N-TriplesTests/TESTS.tar.gz | tar xvz - -C ntriples/testdata/suite
	rm nquads/testdata/suite/*   && curl -s -L https://www.w3.org/2013/N-QuadsTests/TESTS.tar.gz   | tar xvz - -C nquads/testdata/suite
	rm turtle/testdata/suite/*   && curl -s -L https://www.w3.org/2013/TurtleTests/TESTS.tar.gz    | tar xvz - -C turtle/testdata/suite
	mv turtle/testdata/suite/TurtleTests/* turtle/testdata/suite && rmdir turtle/testdata/suite/TurtleTests # move files up one level
	rm trig/testdata/suite/*     && curl -s -L https://www.w3.org/2013/TrigTests/TESTS.tar.gz	   | tar xvz - -C trig/testdata/suite
