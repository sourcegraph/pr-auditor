#! /usr/bin/env bash

set -eu

echo "--- Running 'pr-auditor'"
go run . \
  -github.payload-path="$GITHUB_EVENT_PATH" \
  -github.token="$GITHUB_TOKEN" \
  -github.run-url="$GITHUB_RUN_URL"
