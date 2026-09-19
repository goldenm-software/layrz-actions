# Dart/Flutter Checks

Run Flutter analyze and tests with coverage reporting.

## Usage

```yaml
- name: Run Dart/Flutter checks
  uses: goldenm-software/layrz-actions/.github/actions/check-dart@v1
  with:
    working-directory: '.'
    flutter-version: '3.38.8'
    run-checks: 'true'
    run-tests: 'true'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `working-directory` | Directory to run Flutter commands in | No | `.` |
| `flutter-version` | Flutter version to use | No | `3.38.8` |
| `run-checks` | Whether to run Flutter analyze | No | `true` |
| `run-tests` | Whether to run tests | No | `true` |
| `upload-artifact` | Whether to upload coverage artifact | No | `true` |
| `artifact-name` | Name for the coverage artifact | No | `flutter-coverage` |
| `pub-dry-run` | Whether to run `flutter pub publish --dry-run` | No | `false` |
| `ignore-publish-advisories` | When `true`, ignore three known-benign publish advisories (pre-release dependency, dartbug.com analyzer quirk, git-source); when `false`, any advisory fails the dry-run (strict). Only takes effect when `pub-dry-run` is `true` | No | `false` |
| `ignore-missing-metadata-files` | When `true`, ignore the blocking "Please add a README.md" and "Please add a CHANGELOG.md" issues from the publish dry-run. Use for internal monorepo packages that are never published and intentionally omit those files. Only takes effect when `pub-dry-run` is `true`, and is independent of `ignore-publish-advisories` | No | `false` |

## Outputs

| Output | Description |
|--------|-------------|
| `coverage-path` | Path to the coverage/lcov.info file |
| `test-results-path` | Path to the test-results.json file |
| `artifact-id` | ID of the uploaded artifact |
| `artifact-url` | URL of the uploaded artifact |

## What It Does

1. **Setup Flutter**: Configures Flutter SDK with caching based on pubspec.yaml
2. **Install dependencies**: Runs `flutter pub get` to install packages
3. **Run checks**: Executes `flutter analyze` for static analysis
4. **Run tests**: Runs `flutter test --machine --coverage` to generate coverage data
5. **Publish test results**: Uses test-reporter to display test results in the UI
6. **Upload coverage artifact**: Uploads the lcov.info file as an artifact (retention: 1 day)
7. **Run flutter pub get with --dry-run** *(optional, `pub-dry-run: true`)*: Runs `flutter pub publish --dry-run` and evaluates its output; `ignore-publish-advisories` and `ignore-missing-metadata-files` each relax a different, independent set of blocking issues

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
