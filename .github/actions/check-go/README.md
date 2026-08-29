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
| `golangci-lint-version` | Version of golangci-lint to use. `auto` (default) resolves it via the precedence below; pass an explicit version (e.g. `v2.11.4`) to pin it | No | `auto` |
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
2. **Run linter**: Resolves the golangci-lint version, then executes golangci-lint for code quality checks
3. **Run tests**: Runs `go test -v -coverprofile=coverage.out ./...` to generate coverage
4. **Generate summary**: Uses `go tool cover` to create a human-readable coverage summary
5. **Upload coverage artifact**: Uploads both coverage.out and coverage-summary.txt (retention: 1 day)

### golangci-lint version resolution

When `golangci-lint-version` is left at its default (`auto`), the "Resolve golangci-lint version"
step picks a version using this precedence, highest first:

1. **Explicit input** — if `golangci-lint-version` is set to anything other than `auto`, it is
   honoured verbatim and nothing below is consulted.
2. **`<working-directory>/.golangci-lint-version`** — if this file exists in the consuming repo,
   its contents are used (passed to `golangci-lint-action` via `version-file`).
3. **`<working-directory>/.tool-versions`** — if `.golangci-lint-version` is absent but this file
   exists *and* contains a line starting with `golangci-lint` (the asdf/mise format), it is used
   the same way. A `.tool-versions` file without a `golangci-lint` entry is ignored.
4. **Hardcoded fallback** — if none of the above apply, the action falls back to a previously
   hardcoded default version (currently `v2.11`).

A pre-flight compatibility check also runs regardless of which source won: golangci-lint's own
`go.mod` is pinned by upstream policy to (latest Go minor - 1), so if the consuming repo's `go.mod`
targets a newer Go language version than golangci-lint currently supports, a warning is logged
explaining that the "Run linter" step will likely fail and how to fix it (lower the `go` directive
in `go.mod` — this does not affect the Go toolchain used to build/run the code, only the minimum
language version golangci-lint's config parser expects).

**Note:** `version-file` only works with `golangci-lint-action`'s default `install-mode: binary`
(cached), which this action uses.

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
