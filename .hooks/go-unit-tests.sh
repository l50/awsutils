#!/bin/bash

set -x

TESTS_TO_RUN=$1
RETURN_CODE=0

if [[ -z "${TESTS_TO_RUN}" ]]; then
  echo "No tests input"
  echo "Example - Run all tests: bash go-unit-tests.sh all"
  echo "Example - Run ec2 tests: bash go-unit-tests.sh ec2"
  exit 1
fi

# If we are in an action, run the coverage test.
if [[ "${GITHUB_ACTIONS}" == "true" ]]; then
  TESTS_TO_RUN='coverage'
fi

if [[ "${TESTS_TO_RUN}" == 'coverage' ]]; then
  go test -v -race -failfast \
    -tags=integration -coverprofile=coverage-all.out ./...
  RETURN_CODE=$?
elif [[ "${TESTS_TO_RUN}" == 'all' ]]; then
  go test -v -count=1 -race ./...
  RETURN_CODE=$?
elif [[ "${TESTS_TO_RUN}" == 'short' ]]; then
  go test -v -count=1 -short -race ./...
  RETURN_CODE=$?
else
  go test -v -count=1 -race "./.../${TESTS_TO_RUN}"
  RETURN_CODE=$?
fi

if [[ "${RETURN_CODE}" -ne 0 ]]; then
  echo "unit tests failed"
  exit 1
fi
