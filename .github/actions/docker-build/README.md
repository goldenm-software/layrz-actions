# Build Docker and Push to GitHub Container Registry

Build Docker image for a specific architecture and push to container registry.

## Usage

```yaml
- name: Build Docker image
  uses: goldenm-software/layrz-actions/.github/actions/docker-build@v1
  with:
    platform: linux/amd64
    arch: amd64
    registry: ghcr.io
    username: ${{ github.repository_owner }}
    password: ${{ secrets.GITHUB_TOKEN }}
    build-args: |
      NODE_ENV=production
      APP_VERSION=${{ github.ref_name }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `platform` | Target platform (e.g., linux/amd64, linux/arm64) | Yes | - |
| `arch` | Target architecture for image tags (e.g., amd64, arm64) | Yes | - |
| `registry` | Container registry URL (e.g., ghcr.io, 123456.dkr.ecr.us-east-1.amazonaws.com) | Yes | - |
| `username` | Registry username | Yes | - |
| `password` | Registry password | Yes | - |
| `build-args` | Build arguments (newline-separated) | No | `` |

## Outputs

This action does not have outputs. Images are pushed directly to the registry.

## What It Does

1. **Setup Buildx**: Configures Docker Buildx for multi-platform builds
2. **Login to registry**: Authenticates with the specified container registry
3. **Prepare build args**: Combines default and custom build arguments
   - Default: `BUILDKIT_INLINE_CACHE=1`
   - Custom: Any additional build args provided
4. **Build and push**: Builds the Docker image and pushes it with tags:
   - `sha-{github.sha}-{arch}`: SHA-specific tag
   - `latest-{arch}`: Latest tag for the architecture
5. **Layer caching**: Uses registry cache for faster subsequent builds

## Image Tags

Images are tagged with both SHA and latest tags:
- `{registry}/{repository}:sha-{sha}-{arch}`
- `{registry}/{repository}:latest-{arch}`

**Note:** The `arch` parameter is used in tags instead of `platform` because Docker tags cannot contain `/` characters.

## Full Examples

See the [main repository documentation](https://github.com/goldenm-software/layrz-actions) for complete workflow examples.
