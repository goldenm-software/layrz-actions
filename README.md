# Layrz Actions

Reusable GitHub Actions composite actions for standardized CI/CD workflows across Layrz projects.

## Quick Start

```yaml
jobs:
  lint-python:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: goldenm-software/layrz-actions/.github/actions/check-python@v1
        with:
          run-tests: true
          coverage-library: "myapp"
```

## Benefits

- **Flat check names** - Works seamlessly with branch protection rules
- **Maximum flexibility** - Add custom steps before or after any action
- **No copy-paste** - Update workflows once, all repos benefit
- **Consistency** - Same CI/CD process across all projects
- **Artifact outputs** - All check actions output artifact IDs, URLs, and paths

## Available Actions

### Code Quality & Testing

| Action | Description | Documentation |
|--------|-------------|---------------|
| [check-python](/.github/actions/check-python) | Python linting (Ruff), type checking, tests with UV | [README](/.github/actions/check-python/README.md) |
| [check-dart](/.github/actions/check-dart) | Dart/Flutter analyze and tests with coverage | [README](/.github/actions/check-dart/README.md) |
| [check-go](/.github/actions/check-go) | Go linting (golangci-lint) and tests | [README](/.github/actions/check-go/README.md) |
| [check-kotlin](/.github/actions/check-kotlin) | Kotlin linting (ktlint) and Gradle unit tests | [README](/.github/actions/check-kotlin/README.md) |
| [check-swift](/.github/actions/check-swift) | Swift linting (SwiftLint) and XCTest tests | [README](/.github/actions/check-swift/README.md) |

### Coverage & Changelog

| Action | Description | Documentation |
|--------|-------------|---------------|
| [coverage-comment](/.github/actions/coverage-comment) | Unified coverage reporting for Python, Go, Dart | [README](/.github/actions/coverage-comment/README.md) |
| [changelog](/.github/actions/changelog) | Auto-generate changelog for PRs and releases | [README](/.github/actions/changelog/README.md) |

### Docker

| Action | Description | Documentation |
|--------|-------------|---------------|
| [docker-build](/.github/actions/docker-build) | Platform-specific Docker builds with custom build-args | [README](/.github/actions/docker-build/README.md) |

## Usage Examples

### Python Project with Coverage

```yaml
name: Python Checks

on:
  pull_request:
    branches: [main]

permissions:
  contents: read
  pull-requests: write

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - name: Run Python checks
        id: python-checks
        uses: goldenm-software/layrz-actions/.github/actions/check-python@v1
        with:
          run-tests: true
          coverage-library: "myapp"
          upload-artifact: true

      - name: Post coverage comment
        uses: goldenm-software/layrz-actions/.github/actions/coverage-comment@v1
        with:
          enable-python: true
```

### Dart/Flutter Checks

```yaml
name: Dart Checks

on:
  pull_request:
    branches: [main]

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - uses: goldenm-software/layrz-actions/.github/actions/check-dart@v1
        with:
          flutter-version: "3.38.8"
          run-tests: true
```

### Go Checks

```yaml
name: Go Checks

on:
  pull_request:
    branches: [main]

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - uses: goldenm-software/layrz-actions/.github/actions/check-go@v1
        with:
          run-tests: true
```

### Kotlin/Android Checks

**Basic usage (standalone Gradle project):**

```yaml
name: Kotlin Checks

on:
  pull_request:
    branches: [main]

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - uses: goldenm-software/layrz-actions/.github/actions/check-kotlin@v1
        with:
          gradle-task: testDebugUnitTest
          run-lint: true
          run-tests: true
```

**For Flutter plugin with example app (Android module built via example):**

```yaml
- uses: goldenm-software/layrz-actions/.github/actions/check-kotlin@v1
  with:
    working-directory: example/android
    gradle-task: testDebugUnitTest
    test-results-path: ../build/layrz_push/test-results
    run-lint: true
    run-tests: true
```

### Swift/iOS Checks

Note: This action requires `runs-on: macos-latest` or similar macOS runner.

**Recommended: Lint only hand-written sources (excludes Pods, generated code):**

```yaml
name: Swift Checks

on:
  pull_request:
    branches: [main]

jobs:
  lint-and-test:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v6

      - uses: goldenm-software/layrz-actions/.github/actions/check-swift@v1
        with:
          working-directory: example/ios
          lint-directory: ../layrz_push/Sources
          run-lint: true
          run-tests: true
```

**Full directory lint (uses repo .swiftlint.yml if present, otherwise excludes *.g.swift, Pods):**

```yaml
- uses: goldenm-software/layrz-actions/.github/actions/check-swift@v1
  with:
    working-directory: example/ios
    lint-directory: .
    run-lint: true
    run-tests: true
```

### Docker Build (Multi-Architecture)

```yaml
name: Docker Build

on:
  push:
    branches: [main]

permissions:
  contents: read
  packages: write

jobs:
  build-amd64:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - uses: goldenm-software/layrz-actions/.github/actions/docker-build@v1
        with:
          platform: linux/amd64
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}

  build-arm64:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - uses: goldenm-software/layrz-actions/.github/actions/docker-build@v1
        with:
          platform: linux/arm64
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}
```

### Custom Steps Before/After

```yaml
jobs:
  lint-with-custom-steps:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - name: Setup custom environment
        run: echo "Setting up environment..."

      - uses: goldenm-software/layrz-actions/.github/actions/check-python@v1
        with:
          run-tests: true

      - name: Additional validation
        run: echo "Running additional checks..."
```

## Versioning

All examples use `@v1` which automatically gets the latest stable v1.x.x release:

```yaml
# Recommended: automatic updates within v1.x.x
- uses: goldenm-software/layrz-actions/.github/actions/check-python@v1

# Pin to exact version (no updates)
- uses: goldenm-software/layrz-actions/.github/actions/check-python@v1.0.1

# Not recommended for production
- uses: goldenm-software/layrz-actions/.github/actions/check-python@main
```

- `@v1` is a movable tag pointing to the latest v1.x.x release
- Patch releases automatically update the `v1` tag
- Major releases (v2.0.0) may include breaking changes

## Maintainer Guide

### Creating New Releases

```bash
# Show current version
make current-version

# Patch release (v1.0.1 → v1.0.2)
make push

# Minor release (v1.0.2 → v1.1.0)
make push-minor

# Major release (v1.x.x → v2.0.0)
make push-major
```

See [scripts/README.md](/scripts/README.md) for details.

## Requirements

- GitHub Actions enabled in your repository
- Appropriate permissions (`contents: read`, `pull-requests: write`, `packages: write`)
- Docker BuildKit enabled for docker-build (default in GitHub Actions)

## Contributing

See [CLAUDE.md](/CLAUDE.md) for development guidelines and conventions.

## License

Maintained by [Golden M Software](https://github.com/goldenm-software) for use in Layrz projects.
