# Generate Changelog

Generate and optionally post changelog comments for PRs or releases.

## Usage

```yaml
- name: Generate changelog
  uses: goldenm-software/layrz-actions/.github/actions/changelog@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    post-comment: 'true'
    exclude-types: 'chore,style'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `exclude-dependabot` | Skip action for dependabot PRs | No | `true` |
| `exclude-types` | Commit types to exclude from changelog (comma-separated) | No | `chore,style` |
| `comment-header` | Header text for the PR comment | No | `## 📋 Changelog Summary` |
| `from-ref` | Starting reference (commit/tag) for changelog. If empty, auto-detects based on event type | No | `` |
| `to-ref` | Ending reference (commit/tag) for changelog. If empty, auto-detects based on event type | No | `` |
| `github-token` | GitHub token for posting comments (only needed if posting to PR) | No | `` |
| `post-comment` | Whether to post comment to PR (requires github-token) | No | `true` |

## Outputs

| Output | Description |
|--------|-------------|
| `changelog` | Generated changelog text |

## What It Does

1. **Check for dependabot**: Optionally skips execution for dependabot PRs
2. **Determine refs**: Auto-detects git references based on event type (PR, tag push, or commit)
3. **Generate changelog**: Uses requarks/changelog-action to create structured changelog
4. **Fallback generation**: Falls back to simple commit message list if structured changelog fails
5. **Format changelog**: Adds header, footer, and formatting for readability
6. **Create/update PR comment**: Posts or updates a changelog comment on pull requests

## Reference Detection

- **Pull Request**: Compares head SHA with base SHA
- **Tag Push**: Compares current tag with previous tag
- **Default**: Compares current commit with previous commit
- **Custom**: Use `from-ref` and `to-ref` for manual control

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
