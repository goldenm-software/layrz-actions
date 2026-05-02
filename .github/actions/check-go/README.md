# Go Checks

Run Go linting and tests with coverage reporting.

## Usage

```yaml
- name: Run Go checks
  uses: goldenm-software/layrz-actions/.github/actions/check-go@v1
  with:
    working-directory: '.'
    go-version: '1.25.5'
    run-tests: 'true'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `working-directory` | Directory to run Go commands in | No | `.` |
| `go-version` | Go version to use | No | `1.25.5` |
| `go-cache-dependency-path` | Dependency path for Go cache | No | `go.sum` |
| `run-tests` | Whether to run tests | No | `true` |
| `upload-artifact` | Whether to upload coverage artifact | No | `true` |
| `artifact-name` | Name for the coverage artifact | No | `go-coverage` |

## Outputs

| Output | Description |
|--------|-------------|
| `coverage-out-path` | Path to the coverage.out file |
| `coverage-summary-path` | Path to the coverage-summary.txt file |
| `artifact-id` | ID of the uploaded artifact |
| `artifact-url` | URL of the uploaded artifact |

## What It Does

1. **Setup Go**: Configures Go environment with version-specific caching
2. **Run linter**: Executes golangci-lint (v2.8) for code quality checks
3. **Run tests**: Runs `go test -v -coverprofile=coverage.out ./...` to generate coverage
4. **Generate summary**: Uses `go tool cover` to create a human-readable coverage summary
5. **Upload coverage artifact**: Uploads both coverage.out and coverage-summary.txt (retention: 1 day)

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
