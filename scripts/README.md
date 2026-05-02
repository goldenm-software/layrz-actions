# Version Management Scripts

These scripts automate the versioning and release process for GitHub Actions workflows.

## Scripts

### `current-version.sh`
Shows the current version from git tags.

```bash
./scripts/current-version.sh
# Output: Current version: v1.0.1
```

### `version-push.sh` (Patch)
Creates a new patch version and pushes to GitHub.

**Use for:** Bug fixes, small improvements

**What it does:**
1. Gets current version (e.g., `v1.0.1`)
2. Increments patch number → `v1.0.2`
3. Creates tag `v1.0.2`
4. Moves major version tag `v1` to point to `v1.0.2`
5. Pushes both tags to GitHub

**Usage:**
```bash
./scripts/version-push.sh
# Or via Makefile:
make push
```

### `version-push-minor.sh` (Minor)
Creates a new minor version and pushes to GitHub.

**Use for:** New features, backwards-compatible changes

**What it does:**
1. Gets current version (e.g., `v1.0.2`)
2. Increments minor, resets patch → `v1.1.0`
3. Creates tag `v1.1.0`
4. Moves major version tag `v1` to point to `v1.1.0`
5. Pushes both tags to GitHub

**Usage:**
```bash
./scripts/version-push-minor.sh
# Or via Makefile:
make push-minor
```

### `version-push-major.sh` (Major)
Creates a new major version and pushes to GitHub.

**Use for:** Breaking changes

**What it does:**
1. Gets current version (e.g., `v1.2.3`)
2. Increments major, resets minor/patch → `v2.0.0`
3. Creates tag `v2.0.0`
4. Creates major version tag `v2` pointing to `v2.0.0`
5. Pushes both tags to GitHub

**Usage:**
```bash
./scripts/version-push-major.sh
# Or via Makefile:
make push-major
```

## Semantic Versioning

We follow [Semantic Versioning](https://semver.org/) (SemVer): `MAJOR.MINOR.PATCH`

- **PATCH** (v1.0.1 → v1.0.2): Bug fixes, small improvements
- **MINOR** (v1.0.2 → v1.1.0): New features, backwards-compatible
- **MAJOR** (v1.2.3 → v2.0.0): Breaking changes

**Quick Reference:**
```bash
make push         # Patch: Bug fixes (v1.0.1 → v1.0.2)
make push-minor   # Minor: New features (v1.0.2 → v1.1.0)
make push-major   # Major: Breaking changes (v1.2.3 → v2.0.0)
```

## How Versioning Works

### Tag Strategy

Each release creates two tags:

1. **Specific version tag** (e.g., `v1.0.2`)
   - Never changes
   - Used for pinning to exact versions

2. **Major version tag** (e.g., `v1`)
   - Moves to latest patch/minor within that major version
   - Used for automatic updates

### Example Flow

```bash
# Release v1.0.1
git tag v1.0.1
git tag v1  # points to v1.0.1

# Release v1.0.2
git tag v1.0.2
git tag -f v1  # now points to v1.0.2

# Users referencing @v1 automatically get v1.0.2
```

## Usage in Workflows

Users can reference workflows using:

```yaml
# Automatic updates (recommended)
- uses: goldenm-software/layrz-actions/.github/actions/python-checks@v1

# Pinned to specific version
- uses: goldenm-software/layrz-actions/.github/actions/python-checks@v1.0.2
```

## Requirements

- Git repository with at least one commit
- Permission to push tags to origin
- Bash shell
