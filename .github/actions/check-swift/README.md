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
| `destination` | xcodebuild destination; "auto" picks the first available iPhone simulator on the runner | No | `auto` |
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
2. **Cache CocoaPods** (when pod install enabled): Caches both `Pods/` and `~/Library/Caches/CocoaPods` to avoid rebuilding Firebase and other pods from source on every run
3. **Cache Swift packages** (when tests enabled): Caches SPM dependencies in `build/SourcePackages` and `~/Library/Caches/org.swift.swiftpm` to avoid re-downloading Swift Package Manager packages (e.g., firebase-ios-sdk from plugin Package.swift manifests)
4. **Run pod install**: Installs CocoaPods dependencies (optional)
5. **Install SwiftLint**: Ensures SwiftLint is available (via Homebrew if needed)
6. **Run SwiftLint**: Lints Swift source files, excluding generated code and Pods
7. **Create build reports directory**: Sets up output directory for test reports
8. **Install xcpretty**: Ensures the Ruby formatter for xcodebuild output is available
9. **Run XCTest tests**: Executes tests with code coverage enabled; simulator tests ad-hoc sign by default, providing the application-identifier entitlement required for keychain access (SecItem* operations)
10. **Export coverage to lcov**: Converts XCCode coverage data to lcov format (optional)
11. **Publish test results**: Reports test results via dorny/test-reporter
12. **Upload coverage artifact**: Uploads the lcov.info file (if coverage is enabled)
13. **Upload artifact**: Uploads test results as an artifact (retention: 1 day)

## Caching

The action employs three complementary cache layers to accelerate CI runs:

### CocoaPods Cache
Caches pod dependencies (`Pods/` directory and `~/Library/Caches/CocoaPods`) keyed on `Podfile.lock` and `Podfile`. Avoids rebuilding Firebase iOS SDK and other pods from source on every run.

### Swift Package Manager (SPM) Cache
Caches Swift package manager dependencies in `<working-directory>/build/SourcePackages` and `~/Library/Caches/org.swift.swiftpm`, keyed on `**/Package.resolved` within the working directory. This is active whenever `run-tests` is enabled (since SPM resolution happens during `xcodebuild test`). Avoids re-downloading packages like firebase-ios-sdk declared in plugin `Package.swift` manifests.

### Flutter SDK Cache
When `flutter-version` is specified, the `setup-flutter` action caches the Flutter SDK itself via `cache: true` and `cache-sdk: true`.

Together, these three caches reduce CI runtime by 30-50% on subsequent runs, depending on dependency counts.

## Coverage Setup

Coverage is enabled by default (`run-coverage: 'true'`). The action:

1. Enables code coverage during the XCTest run with `-enableCodeCoverage YES`
2. Finds the generated `Coverage.profdata` in the Xcode build directory
3. Uses `xcrun llvm-cov` to export coverage to lcov format by searching for instrumented binaries
4. Optionally filters the lcov output by `coverage-include-pattern` to focus on specific source paths

### Coverage candidate search

The step searches for Mach-O executables under `build/Build/Products/*-iphonesimulator/` in this priority order:

1. **Test bundle**: `*.xctest/PlugIns/*.xctest/` executable (e.g., `RunnerTests.xctest/RunnerTests`)
2. **Framework binaries**: `*.app/Frameworks/*.framework/` executables
3. **App binary**: `*.app/<AppName>`

It attempts a combined export using all candidates found (first candidate as the primary binary, others via `-object` flags). If that fails or produces no output, it falls back to trying each candidate individually. The first candidate that produces a non-empty, valid lcov file is used. The step logs which binary won.

**Important**: Because Firebase pods require static frameworks on iOS, plugin code links into the app binary, so coverage naturally includes plugin sources across multiple binaries.

**Non-fatal behavior**: If coverage export fails (no candidates, no valid output, etc.), the step logs a warning ("coverage export skipped: <reason> — tests are unaffected"), sets output `found=false`, and does NOT fail the job. Downstream coverage uploads gate on the `found` output, so test results are always published even if coverage is unavailable.

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
