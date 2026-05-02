#!/bin/bash
# Show current version
# Usage: ./scripts/current-version.sh

# Get current version (only semver tags like v1.2.3, ignore v1, v2, etc.)
CURRENT_VERSION=$(git tag -l 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | head -n1 || echo "No tags found")
echo "Current version: $CURRENT_VERSION"
