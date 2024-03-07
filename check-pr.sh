#! /usr/bin/env bash
set -eu

echo "--- Running 'pr-auditor'"
go run . \
  -github.payload-path="$GITHUB_EVENT_PATH" \
  -github.token="$GITHUB_TOKEN" \
  -github.run-url="$GITHUB_RUN_URL" \
  -skip-check-test-plan="${SKIP_CHECK_TEST_PLAN:-False}" \
  -skip-check-review="${SKIP_CHECK_REVIEWS:-False}" \
  -skip-check-review-for-users="${SKIP_CHECK_REVIEWS_FOR_USERS:-""}"
