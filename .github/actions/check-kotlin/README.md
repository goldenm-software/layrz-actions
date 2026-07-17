# Kotlin/Android Checks

Run Kotlin linting (ktlint) and Gradle unit tests for Android with optional JaCoCo coverage reporting.

## Usage

```yaml
- name: Run Kotlin checks
  uses: goldenm-software/layrz-actions/.github/actions/check-kotlin@v1
  with:
    working-directory: 'android'
    gradle-task: 'testDebugUnitTest'
    run-lint: 'true'
    run-tests: 'true'
    coverage-path: 'build/reports/jacoco/testDebugUnitTest/jacoco.xml'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `working-directory` | Directory containing the Gradle project | No | `.` |
| `gradle-task` | Gradle task to run for tests | No | `testDebugUnitTest` |
| `java-version` | Java version to use (temurin distribution) | No | `17` |
| `ktlint-version` | Version of ktlint CLI to use | No | `1.4.1` |
| `flutter-version` | Flutter version to use; when non-empty, set up Flutter first | No | `` |
| `flutter-directory` | Directory to run flutter pub get in | No | `.` |
| `run-lint` | Whether to run ktlint linting | No | `true` |
| `run-tests` | Whether to run Gradle tests | No | `true` |
| `upload-artifact` | Whether to upload test results artifact | No | `true` |
| `artifact-name` | Name for the test results artifact | No | `kotlin-test-results` |
| `test-results-path` | Relative path to test results directory (from working-directory) | No | `build/test-results` |
| `coverage-path` | Path to the JaCoCo XML report, relative to working-directory; empty disables coverage upload | No | `` |
| `coverage-artifact-name` | Name for the coverage artifact | No | `kotlin-coverage` |

## Outputs

| Output | Description |
|--------|-------------|
| `test-results-path` | Path to the Gradle test results XML directory |
| `coverage-path` | Path to the JaCoCo XML coverage file |
| `artifact-id` | ID of the uploaded artifact |
| `artifact-url` | URL of the uploaded artifact |

## What It Does

1. **Setup Flutter** (optional): Configures Flutter if flutter-version is specified
2. **Setup Java**: Configures the Java environment with the specified version
3. **Download ktlint**: Retrieves the ktlint CLI tool if lint is enabled
4. **Run ktlint**: Lints Kotlin source files, excluding generated files (*.g.kt)
5. **Run Gradle tests**: Executes the specified Gradle test task
6. **Publish test results**: Reports test results via dorny/test-reporter
7. **Upload coverage artifact**: Uploads the JaCoCo XML report (if coverage-path is specified)
8. **Set outputs**: Provides paths to test results and coverage files

## Coverage Setup

To enable JaCoCo coverage reporting, your Gradle build must be configured to generate a JaCoCo XML report:

```kotlin
// build.gradle.kts
plugins {
    id("jacoco")
}

jacoco {
    toolVersion = "0.8.10"
}

tasks.withType<Test> {
    finalizedBy(tasks.jacocoTestReport)
}

tasks.jacocoTestReport {
    dependsOn(tasks.test)
    reports {
        xml.required.set(true)
        html.required.set(true)
    }
}
```

The action expects the XML report at the path specified in `coverage-path`. Downstream coverage actions (e.g., `coverage-comment`, `coverage-check`) will default to `kotlin-coverage/jacoco.xml` (artifact-name/file-name).

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
