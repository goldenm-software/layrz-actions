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

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
