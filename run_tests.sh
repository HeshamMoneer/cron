#!/bin/bash

# this script runs the available tests, and opens a report with the tests code coverage

pushd "scheduler"
mkdir -p logs
mkdir -p out

go test -coverprofile=out/coverage.out -v

go tool cover -html=out/coverage.out

popd