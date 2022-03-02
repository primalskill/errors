#!/bin/sh

go test -coverprofile ./cover.out
go tool cover -html=./cover.out -o coverage-report.html