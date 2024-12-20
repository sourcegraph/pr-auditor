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
go run . \
  -github.payload-path="$GITHUB_EVENT_PATH" \
  -github.token="$GITHUB_TOKEN" \
  -protected-branch="release" \
  -skip-check-test-plan=true

# run using wrapper script
./check-pr.sh
```

## Opting out of checks

Each check that PR auditor performs can be opted out of a repository level if they are inappropriate for your use cases. Simply set the relevant environment variable in your GitHub Action to a truthy value like `True` or `true`. By default all checks are enabled.

| Environment Variable         | Check Description                                                                                                                                                                                                      |
| ---------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| SKIP_CHECK_TEST_PLAN         | Allows PRs to not include the Test Plan section. Useful for repositories which do not include source code (such as documentation repos).                                                                               |
| SKIP_CHECK_REVIEWS           | Allows PRs to be merged without requiring reviews. Useful for repositories which are entirely automated (such as infrastructure code).                                                                                 |
| SKIP_CHECK_REVIEWS_FOR_USERS | Allows PRs to be merged without requiring reviews for the specified users. Useful for repositories which have a clear owner(s). Format is CSV of GitHub handles. _Note: This has no effect if SKIP_CHECK_REVIEWS=true_ |

## Deployment

`pr-auditor` can be deployed to repositories using the available [batch changes](./batch-changes/README.md).

You can also add it to a single repo by copying `pr-auditor.example.yml` to `.github/workflows/pr-auditor.yml`.

You will also need to add the `sourcegraph-bot-devx` user to the repository as a collaborator with write access.
