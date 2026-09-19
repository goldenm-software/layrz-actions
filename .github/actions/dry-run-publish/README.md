# Dry Run Publish

Run `flutter pub publish --dry-run` and decide pass/fail by inspecting the
reported output, instead of trusting the dry-run's own exit code alone.

## Usage

```yaml
- uses: goldenm-software/layrz-actions/.github/actions/dry-run-publish@v1
  with:
    flutter-version: '3.47.2'
    ignore-publish-advisories: 'true'
    ignore-missing-metadata-files: 'true'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `working-directory` | Directory to run Flutter commands in | No | `.` |
| `flutter-version` | Flutter version to use | No | `3.41.9` |
| `ignore-publish-advisories` | Whether to ignore three known-benign publish advisories | No | `false` |
| `ignore-missing-metadata-files` | Whether to ignore the blocking "Please add a README.md" and "Please add a CHANGELOG.md" issues | No | `false` |

## The `ignore-publish-advisories` flag

`flutter pub publish --dry-run` can report three advisories that are benign
for a package that is git-only (`publish_to: none`) and intentionally depends
on pre-release and/or git-source dependencies:

1. A package depending on a pre-release of another package should itself be
   published as a pre-release.
2. A spurious analyzer quirk in the dry-run's own embedded analyzer, ending
   with "Please report this at dartbug.com." (the package's own
   `flutter analyze`/`dart analyze` reports no issues at the same time).
3. "Don't depend on `<pkg>` from the git source. Use the hosted source
   instead," for a dependency that is deliberately pulled from git.

Set `ignore-publish-advisories: 'true'` for a git-only package that
intentionally depends on pre-release and/or git-source dependencies, so
those three advisory classes are ignored and everything else still fails
the run.

Leave `ignore-publish-advisories: 'false'` (the default) for a hosted package
published to pub.dev, so strict validation applies: any reported issue,
including the three classes above, fails the run.

## The `ignore-missing-metadata-files` flag

`flutter pub publish --dry-run` blocks on two issues when a package is
missing top-level documentation files:

1. "Please add a `README.md` file."
2. "Please add a `CHANGELOG.md` file."

Set `ignore-missing-metadata-files: 'true'` for an internal monorepo package
that is never published and intentionally omits those files, so only those
two blocking issues are ignored and everything else still fails the run.
This flag is independent of `ignore-publish-advisories`: enable either, both,
or neither depending on the package.

Leave `ignore-missing-metadata-files: 'false'` (the default) for a package
that is expected to ship a `README.md` and `CHANGELOG.md`.

## What It Does

1. **Setup Flutter**: configures the Flutter SDK with caching based on
   `pubspec.yaml`.
2. **Install dependencies**: runs `flutter pub get`.
3. **Validate publish (dry-run)**: runs the bundled `dry_run_publish.sh`,
   which invokes `flutter pub publish --dry-run`, parses the reported issue
   blocks, and fails closed on anything not covered by
   `ignore-publish-advisories` or `ignore-missing-metadata-files`.

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
