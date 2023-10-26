SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

CURPATH := ${shell pwd}

# Filter down to run tests only in this path. To run all the tests from all the project folders use ./...
# To run tests only in root use ./.
# Make sure TESTREGEX below is set to an appropiate value too.
TESTPATH := ./...

# Filter down to testing only a portion of the tests by using a regex to match test names. Ex. TestUpsertUser
# To test multiple at once separate with / ex: TestUpsertUser/TestGetUserIDsByToken
# To test all use dot ex. .
TESTREGEX := TestJSONMarshaling

# Running 'make cover' will create a test coverage report html file in the root of the source folder.
COVERAGEREPORT := ${CURPATH}/coverage-report.html

API_APP := main
API_INPUT := ${CURPATH}/cmd/api/${API_APP}.go
API_OUTPUT := /usr/local/bin/${API_APP}

.PHONY: update-deps-latest
update-deps-latest:
	go mod tidy
	go get -u ./...

.PHONY: update-deps
update-deps:
	go mod tidy
	go get ./...

.PHONY: format
format:
	gofmt -w ${CURPATH}

.PHONY: test
test: format
	clear && printf '\e[3J'	
	go clean -testcache
	go test -coverprofile /tmp/cover.out -v -failfast -run ${TESTREGEX} ${TESTPATH}
	go tool cover -html=/tmp/cover.out -o ${COVERAGEREPORT}