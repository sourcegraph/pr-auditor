# pr-auditor [![pr-auditor](https://github.com/sourcegraph/sourcegraph/actions/workflows/pr-auditor.yml/badge.svg)](https://github.com/sourcegraph/sourcegraph/actions/workflows/pr-auditor.yml)

`pr-auditor` is a tool designed to operate on some [GitHub Actions pull request events](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request) in order to check for SOC2 compliance.
Owned by the [DevX team](https://handbook.sourcegraph.com/departments/product-engineering/engineering/enablement/dev-experience).

Learn more: [Testing principles and guidelines](https://docs.sourcegraph.com/dev/background-information/testing_principles)

## Usage

This action is primarily designed to run on GitHub Actions, and leverages the [pull request event payloads](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request) extensively.

The optional `-protected-branch` flag defines a base branch that always opens a PR audit issue to track all pull requests made to it.

```sh
GITHUB_EVENT_PATH="/path/to/json/payload.json"
GITHUB_TOKEN="personal-access-token"

# run directly
go run ./dev/pr-auditor/ check \
  -github.payload-path="$GITHUB_EVENT_PATH" \
  -github.token="$GITHUB_TOKEN" \
  -protected-branch="release"

# run using wrapper script
./dev/buildchecker/check-pr.sh
```

## Opting out of checks

Each check that PR auditor performs can be opted out of a repository level if they are inappropriate for your use cases. Simply set the relevant environment variable in your GitHub Action to a falsey value like `False` or `false`. By default all checks are enabled.

| Environment Variable | Check Description                                                                                                                        |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| SKIP_CHECK_TEST_PLAN | Allows PRs to not include the Test Plan section. Useful for repositories which do not include source code (such as documentation repos). |
| SKIP_CHECK_REVIEWS   | Allows PRs to be merged without requiring reviews. Useful for repositories which are entirely automated (such as infrastructure code).   |

## Deployment

`pr-auditor` can be deployed to repositories using the available [batch changes](./batch-changes/README.md).
