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
2. **Ensure golangci-lint version file**: Guarantees `<working-directory>/.golangci-lint-version`
   exists before linting
3. **Run linter**: Executes golangci-lint for code quality checks
4. **Run tests**: Runs `go test -v -coverprofile=coverage.out ./...` to generate coverage
5. **Generate summary**: Uses `go tool cover` to create a human-readable coverage summary
6. **Upload coverage artifact**: Uploads both coverage.out and coverage-summary.txt (retention: 1 day)

### golangci-lint version

The linter step always runs via `version-file: <working-directory>/.golangci-lint-version`. The
"Ensure golangci-lint version file" step guarantees that file exists beforehand:

- If the consuming repo already has `<working-directory>/.golangci-lint-version`, it is used
  untouched.
- Otherwise, the action creates it in the runner's checkout containing a hardcoded fallback
  version (currently `v2.11`). This file is ephemeral — it lives only in the CI checkout for the
  duration of the job and is never committed back to the repo.

**Note:** `version-file` only works with `golangci-lint-action`'s default `install-mode: binary`
(cached), which this action uses.

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
