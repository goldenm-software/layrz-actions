# Coverage Report Comment

Parse coverage files and post a unified coverage report comment to PR.

## Usage

```yaml
- name: Post coverage report
  uses: goldenm-software/layrz-actions/.github/actions/coverage-comment@v1
  with:
    enable-python: 'true'
    enable-go: 'true'
    enable-dart: 'false'
    github-token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `enable-python` | Enable Python coverage reporting | No | `false` |
| `enable-go` | Enable Go coverage reporting | No | `false` |
| `enable-dart` | Enable Dart/Flutter coverage reporting | No | `false` |
| `python-coverage-path` | Path to Python coverage.json file (must be downloaded first) | No | `python-coverage/coverage.json` |
| `go-coverage-path` | Path to Go coverage-summary.txt file (must be downloaded first) | No | `go-coverage/coverage-summary.txt` |
| `dart-coverage-path` | Path to Dart lcov.info file (must be downloaded first) | No | `flutter-coverage/lcov.info` |
| `github-token` | GitHub token for posting comments | No | `` |

## Outputs

This action does not have outputs. It posts a comment directly to the PR.

## What It Does

1. **Install lcov**: Installs lcov tool if Dart coverage is enabled
2. **Parse Dart coverage**: Extracts coverage percentage and line counts from lcov.info
3. **Parse coverage files**: Reads and parses coverage data for enabled languages:
   - **Python**: Parses coverage.json for percentage and line metrics
   - **Go**: Extracts total coverage percentage from coverage-summary.txt
   - **Dart**: Uses pre-parsed lcov data
4. **Build unified report**: Creates a single markdown table with all coverage metrics
5. **Post/update comment**: Creates or updates a PR comment with the coverage report

## Prerequisites

Coverage files must be downloaded before using this action:

```yaml
- name: Download coverage artifacts
  uses: actions/download-artifact@v4
  with:
    name: python-coverage
    path: python-coverage/
```

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
