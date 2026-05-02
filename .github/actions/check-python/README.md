# Python Checks

Run Python linting, type checking, and tests with UV and Ruff.

## Usage

```yaml
- name: Run Python checks
  uses: goldenm-software/layrz-actions/.github/actions/check-python@v1
  with:
    working-directory: '.'
    run-type-check: 'true'
    run-tests: 'true'
    coverage-library: 'my_package'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `working-directory` | Directory to run Python commands in | No | `.` |
| `uv-cache-dependency-glob` | Dependency glob pattern for UV cache | No | `uv.lock` |
| `run-type-check` | Whether to run type checking | No | `true` |
| `run-tests` | Whether to run tests with coverage | No | `false` |
| `coverage-library` | Library to measure coverage for | No | `` |
| `upload-artifact` | Whether to upload coverage artifact | No | `true` |
| `artifact-name` | Name for the coverage artifact | No | `python-coverage` |

## Outputs

| Output | Description |
|--------|-------------|
| `coverage-path` | Path to the coverage.json file |
| `artifact-id` | ID of the uploaded artifact |
| `artifact-url` | URL of the uploaded artifact |

## What It Does

1. **Setup UV**: Configures the UV package manager with Python caching enabled
2. **Install dependencies**: Runs `uv sync --frozen` to install project dependencies
3. **Run lint checks**: Executes `uv run ruff check` for code linting
4. **Run type checks**: Optionally runs `uv run ty check` for type validation
5. **Run tests with coverage**: Optionally executes pytest with coverage reporting
6. **Upload coverage artifact**: Uploads the coverage.json file as an artifact (retention: 1 day)

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
