# Swift/iOS Checks

Run SwiftLint linting and XCTest tests for iOS with optional lcov coverage reporting.

## Usage

```yaml
- name: Run Swift checks
  uses: goldenm-software/layrz-actions/.github/actions/check-swift@v1
  with:
    working-directory: 'example/ios'
    workspace: 'Runner.xcworkspace'
    scheme: 'Runner'
    run-lint: 'true'
    run-tests: 'true'
    run-coverage: 'true'
    coverage-include-pattern: 'layrz_push'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `working-directory` | Directory containing the iOS workspace (e.g., example/ios) | No | `example/ios` |
| `workspace` | Name of the Xcode workspace | No | `Runner.xcworkspace` |
| `scheme` | Xcode scheme to build and test | No | `Runner` |
| `destination` | xcodebuild destination specifier | No | `platform=iOS Simulator,name=iPhone 16` |
| `flutter-version` | Flutter version to use; when non-empty, set up Flutter and run flutter build ios --config-only | No | `` |
| `flutter-directory` | Directory to run Flutter commands in | No | `.` |
| `run-pod-install` | Whether to run pod install (cocoapods is preinstalled on macOS) | No | `true` |
| `lint-directory` | Directory path to lint with swiftlint (e.g., ios/layrz_push/Sources for hand-written code) | No | `.` |
| `swiftlint-version` | Version of SwiftLint to use (e.g., 0.57.0); downloaded from official GitHub releases | No | `0.57.0` |
| `run-lint` | Whether to run SwiftLint | No | `true` |
| `run-tests` | Whether to run XCTest tests | No | `true` |
| `upload-artifact` | Whether to upload test results artifact | No | `true` |
| `artifact-name` | Name for the test results artifact | No | `swift-test-results` |
| `run-coverage` | Whether to export coverage to lcov | No | `true` |
| `coverage-artifact-name` | Name for the coverage artifact | No | `swift-coverage` |
| `coverage-include-pattern` | Only source paths matching this substring are kept in the lcov output; empty keeps everything | No | `` |

## Outputs

| Output | Description |
|--------|-------------|
| `test-results-path` | Path to the junit test results XML file |
| `coverage-path` | Path to the lcov coverage file |
| `artifact-id` | ID of the uploaded artifact |
| `artifact-url` | URL of the uploaded artifact |

## What It Does

1. **Setup Flutter** (optional): Configures Flutter if flutter-version is specified
2. **Run pod install**: Installs CocoaPods dependencies (optional)
3. **Install SwiftLint**: Ensures SwiftLint is available (via Homebrew if needed)
4. **Run SwiftLint**: Lints Swift source files, excluding generated code and Pods
5. **Create build reports directory**: Sets up output directory for test reports
6. **Install xcpretty**: Ensures the Ruby formatter for xcodebuild output is available
7. **Run XCTest tests**: Executes tests with code coverage enabled
8. **Export coverage to lcov**: Converts XCCode coverage data to lcov format (optional)
9. **Publish test results**: Reports test results via dorny/test-reporter
10. **Upload coverage artifact**: Uploads the lcov.info file (if coverage is enabled)
11. **Upload artifact**: Uploads test results as an artifact (retention: 1 day)

## Coverage Setup

Coverage is enabled by default (`run-coverage: 'true'`). The action:

1. Enables code coverage during the XCTest run with `-enableCodeCoverage YES`
2. Finds the generated `Coverage.profdata` in the Xcode build directory
3. Uses `xcrun llvm-cov` to export coverage to lcov format
4. Optionally filters the lcov output by `coverage-include-pattern` to focus on specific source paths

**Important**: Because Firebase pods require static frameworks on iOS, plugin code links into the app binary, so the app binary's coverage naturally includes plugin sources.

### Example: Filter coverage by plugin sources

```yaml
- name: Run Swift checks
  uses: goldenm-software/layrz-actions/.github/actions/check-swift@v1
  with:
    working-directory: 'example/ios'
    coverage-include-pattern: 'layrz_push'
```

This keeps only lcov records whose `SF:` (source file) path contains `layrz_push`.

Downstream coverage actions (e.g., `coverage-comment`, `coverage-check`) will default to `swift-coverage/lcov.info` (artifact-name/file-name).

## SwiftLint Configuration

The action downloads a pinned version of SwiftLint (default 0.57.0) from the official GitHub releases, ensuring version consistency between CI and local tooling. To use a different version, set the `swiftlint-version` input:

```yaml
- name: Run Swift checks
  uses: goldenm-software/layrz-actions/.github/actions/check-swift@v1
  with:
    swiftlint-version: '0.58.0'
```

The action looks for `.swiftlint.yml` in the working directory. If not found, it creates a default configuration that excludes:
- Generated files (`*.g.swift`)
- Pods directory
- Flutter symlinks and dart_tool

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
