# CLAUDE.md - Project Guidelines

This document contains development guidelines, architectural decisions, and conventions for the `layrz-actions` project.

## Project Overview

This repository provides reusable GitHub Actions composite actions for standardized CI/CD workflows across Layrz projects. The primary goal is to eliminate copy-pasting workflow files across repositories while maintaining flexibility and consistency.

**Organization:** `goldenm-software`
**Repository:** `layrz-actions`
**Primary Use Case:** Public reusable actions for CI checks, coverage reporting, changelog generation, and Docker builds

## Maintainer Workflow

**Development Model:** Direct commits to `main` branch (no PRs)
- All development happens directly on the `main` branch
- No pull request workflow is used for this repository
- Changes are committed and pushed immediately
- Branch protection rules are bypassed for maintainers

**Commit Message Format:** Conventional Commits (REQUIRED)

All commits MUST follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>: <description>

[optional body]

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

**Required Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `perf`: Performance improvement
- `refactor`: Code refactoring (no functional changes)
- `test`: Adding or updating tests
- `build`: Build system or dependencies
- `ci`: CI/CD configuration
- `sec`: Security patches (e.g., CVE fixes)
- `lab`: Labs / exploratory work
- `exp`: Experimental features or prototypes
- `deps`: Dependency updates
- `revert`: Reverting previous changes
- `chore`: Maintenance tasks (excluded from changelog)
- `style`: Code style changes (excluded from changelog)

**Auto-Commit Triggers:**
When the following phrases are used, changes should be automatically committed and pushed:
- "wrap it up"
- "commit and push"
- "summarize and commit"
- "done, commit"
- "finalize"

## Architecture Decisions

### 1. Composite Actions Over Reusable Workflows

**Decision:** Use composite actions exclusively, not reusable workflows.

**Rationale:**
- **Flat check names:** Composite actions run as steps, producing check names like `lint-python` instead of nested names like `.github/workflows/checks.yaml / lint-python / lint-python`
- **Branch protection compatibility:** GitHub branch protection rules require exact check name matches. Reusable workflows create nested structures that break this requirement.
- **Flexibility:** Composite actions allow users to add custom steps before/after the action in the same job
- **Simplicity:** Easier to understand and debug since they're just grouped steps

### 2. Artifact Outputs for All Check Actions

**Decision:** All check actions (check-python, check-dart, check-go) output artifact metadata.

**Outputs provided:**
- `artifact-id` - GitHub artifact ID
- `artifact-url` - Direct URL to the artifact
- `coverage-path` - Path to the coverage file (format varies by language)

**Rationale:**
- Maximum flexibility for downstream jobs
- Enables custom post-processing of coverage data
- Allows integration with third-party tools

**Example:**
```yaml
- id: python-checks
  uses: goldenm-software/layrz-actions/.github/actions/check-python@v1

- name: Download coverage
  uses: actions/download-artifact@v4
  with:
    name: ${{ steps.python-checks.outputs.artifact-id }}
```

### 3. Generic Registry Approach

**Decision:** Docker actions use generic registry parameters, not provider-specific names.

**Parameters:**
- `registry` (not `ecr-registry`)
- `registry-user` (not `aws-account-id`)
- `registry-repo` (not `ecr-repository`)

**Rationale:**
- Works with any Docker registry (GHCR, DockerHub, ECR, etc.)
- Reduces lock-in to specific vendors
- Cleaner separation of concerns (build vs. push)

### 4. Platform-Specific Docker Builds

**Decision:** Remove QEMU setup; require separate jobs for each platform.

**Approach:**
```yaml
jobs:
  build-amd64:
    steps:
      - uses: goldenm-software/layrz-actions/.github/actions/docker-build@v1
        with:
          platform: linux/amd64

  build-arm64:
    steps:
      - uses: goldenm-software/layrz-actions/.github/actions/docker-build@v1
        with:
          platform: linux/arm64
```

**Rationale:**
- Faster builds (no emulation overhead)
- Clearer job separation in GitHub Actions UI
- Users control which platforms to build

### 5. Auto-Generated Docker Tags

**Decision:** `docker-build` action auto-generates tags based on SHA and platform.

**Auto-generated tags:**
- `sha-{sha}-{platform}` - For platform-specific images
- `latest-{platform}` - For latest platform-specific images

**Rationale:**
- Prevents tag confusion between platform-specific and manifest tags
- Allows intermediate images to be referenced consistently

## Versioning Strategy

### Movable vs. Fixed Tags

**Movable tags:** `v1`, `v2`, `v3`
**Fixed tags:** `v1.0.0`, `v1.0.1`, `v2.0.0`

**How it works:**
1. Create specific version tag: `v1.0.1`
2. Move major version tag: `v1 → v1.0.1`
3. Users referencing `@v1` automatically get the latest v1.x.x

**Recommendation for users:**
- Use `@v1` for automatic updates (recommended)
- Use `@v1.0.1` for stability-critical projects
- Never use `@main` in production

### Creating Releases

**Automated via Makefile:**

```bash
# Patch release: v1.0.1 → v1.0.2
make push

# Minor release: v1.0.2 → v1.1.0
make push-minor

# Major release: v1.x.x → v2.0.0
make push-major

# Check current version
make current-version
```

**Manual process (if needed):**
```bash
git tag -a v1.0.2 -m "Release v1.0.2"
git tag -af v1 -m "Major version v1 (v1.0.2)"
git push origin v1.0.2
git push origin v1 --force
```

**Important:** Always prompt for confirmation before pushing tags.

## Development Workflow

### Directory Structure

```
.github/
  actions/              # All composite actions
    check-python/
      action.yml
      README.md
    check-dart/
    check-go/
    changelog/
    coverage-comment/
    docker-build/
scripts/                # Versioning automation
  current-version.sh
  version-push.sh
  version-push-minor.sh
  version-push-major.sh
  README.md
examples/               # Example workflows
Makefile                # Release automation
README.md               # Main documentation
CLAUDE.md               # This file
```

### Adding a New Action

1. Create the action directory under `.github/actions/my-action/`
2. Create `action.yml` with `using: composite`
3. Create `README.md` with description, inputs/outputs tables, and usage example using `@v1`
4. Add the action to the table in the main `README.md`
5. Add an example workflow to `examples/`
6. Test in a real repository using `@main` before releasing

### Coding Conventions

**Action Files:**
- Use `.yml` extension (not `.yaml`) for consistency
- Exception: `docker-build/action.yaml` (kept for legacy reasons)
- Always include `description` for inputs and outputs
- Use `required: true/false` explicitly
- Shell must always be `bash`

**Inputs:**
- Use kebab-case: `run-type-check`, not `runTypeCheck`
- Boolean defaults should be strings: `'true'` or `'false'`
- Always provide defaults for optional inputs

**Outputs:**
- Clear, descriptive names: `artifact-id`, not `id`
- Use step IDs for referencing: `${{ steps.step-id.outputs.value }}`

**Bash Scripts:**
- Always use `set -e`
- Quote variables: `"$VARIABLE"`, not `$VARIABLE`
- Use SCREAMING_SNAKE_CASE for variable names

## Language-Specific Conventions

### Python (check-python)
- **Package manager:** UV
- **Linter:** Ruff
- **Coverage format:** `coverage.json`
- **Artifact name:** `python-coverage` (default)

### Dart/Flutter (check-dart)
- **Tool:** Flutter CLI
- **Default version:** `3.38.8`
- **Coverage format:** `lcov.info`
- **Artifact name:** `flutter-coverage` (default)

### Go (check-go)
- **Linter:** golangci-lint
- **Coverage formats:** `coverage.out`, `coverage-summary.txt`
- **Artifact name:** `go-coverage` (default)

### Docker (docker-build)
- **BuildKit:** Always enabled
- **Default build-args:** `BUILDKIT_INLINE_CACHE=1`
- **Tag format:** `sha-{sha}-{platform}`, `latest-{platform}`

## Common Patterns

### Conditional Steps

```yaml
- name: Conditional step
  if: ${{ inputs.run-tests == 'true' }}
  shell: bash
  run: echo "Running tests"
```

### Setting Outputs

```yaml
- name: Set output
  id: set-output
  shell: bash
  run: echo "value=hello" >> $GITHUB_OUTPUT

outputs:
  my-output:
    value: ${{ steps.set-output.outputs.value }}
```

### Uploading Artifacts

```yaml
- name: Upload artifact
  id: upload
  uses: actions/upload-artifact@v4
  with:
    name: my-artifact
    path: path/to/file

outputs:
  artifact-id:
    value: ${{ steps.upload.outputs.artifact-id }}
  artifact-url:
    value: ${{ steps.upload.outputs.artifact-url }}
```

## Troubleshooting

### Nested Check Names
**Problem:** Check appears as `.github/workflows/file.yaml / job / check`
**Solution:** Use composite action, not reusable workflow

### Permission Denied Errors
**Problem:** Action can't write to certain directories
**Solution:** Use `${{ github.workspace }}` or `${{ runner.temp }}`

### Artifact Not Found
**Problem:** Downstream job can't find artifact
**Solution:** Verify artifact name matches exactly and upload succeeded

### Tags Not Updating
**Problem:** Users still see old version when using `@v1`
**Solution:** Force push the movable tag: `git push origin v1 --force`

## Breaking Changes Protocol

1. **Major version bump required** (v1 → v2)
2. **Conventional commit with BREAKING CHANGE footer:**
   ```
   feat: change input parameter name

   BREAKING CHANGE: `registry-user` renamed to `username`

   Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
   ```
3. Update all README files and examples
4. Provide migration guide in release notes

## Key Principles

1. **Flat check names always** - Critical for branch protection
2. **Generic, not specific** - Avoid vendor lock-in
3. **Document everything** - Each action has its own README
4. **Test before releasing** - Use real repositories, not just local tests
5. **Version responsibly** - Major bumps = breaking changes
6. **Artifact outputs** - Provide maximum flexibility
7. **Fail fast** - Use `set -e` in bash scripts

## Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Actions: Composite Actions](https://docs.github.com/en/actions/creating-actions/creating-a-composite-action)
- [Docker BuildKit Documentation](https://docs.docker.com/build/buildkit/)
- [GitHub Actions Artifacts](https://docs.github.com/en/actions/using-workflows/storing-workflow-data-as-artifacts)

## Contact & Support

- **Maintainer:** Golden M Software
- **Issues:** Use GitHub Issues for bug reports and feature requests

---

**Last Updated:** 2026-05-02
**Version:** 1.0.0
