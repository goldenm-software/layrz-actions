#!/bin/bash
# Create a new minor version and push to GitHub
# Usage: ./scripts/version-push-minor.sh

set -e

echo "Creating new minor version..."

# Get current version (only semver tags like v1.2.3, ignore v1, v2, etc.)
CURRENT_VERSION=$(git tag -l 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | head -n1 | sed 's/^v//' || echo "")

if [ -z "$CURRENT_VERSION" ]; then
    echo "Error: No existing tags found. Creating initial version v1.0.0..."
    NEW_VERSION="1.0.0"
    MAJOR_VERSION="1"
else
    echo "Current version: v$CURRENT_VERSION"

    MAJOR=$(echo $CURRENT_VERSION | cut -d. -f1)
    MINOR=$(echo $CURRENT_VERSION | cut -d. -f2)

    NEW_MINOR=$((MINOR + 1))
    NEW_VERSION="$MAJOR.$NEW_MINOR.0"
    MAJOR_VERSION="$MAJOR"
fi

echo "New minor version: v$NEW_VERSION"
echo ""

# Confirm (skip if non-interactive)
if [ -t 0 ]; then
    read -p "Create and push v$NEW_VERSION? [y/N] " -n 1 -r
    echo ""

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Cancelled."
        exit 0
    fi
else
    echo "Non-interactive mode: proceeding automatically."
fi

# Create and push tags
echo "Creating tag v$NEW_VERSION..."
git tag -a "v$NEW_VERSION" -m "Release v$NEW_VERSION"

echo "Moving v$MAJOR_VERSION tag to v$NEW_VERSION..."
git tag -fa "v$MAJOR_VERSION" -m "Update v$MAJOR_VERSION to v$NEW_VERSION"

echo "Pushing tags to origin..."
git push origin "v$NEW_VERSION"
git push origin "v$MAJOR_VERSION" --force

echo ""
echo "✅ Successfully created and pushed v$NEW_VERSION"
echo "✅ Updated v$MAJOR_VERSION tag"
echo ""
echo "Users can now use:"
echo "  @v$MAJOR_VERSION (gets latest v$MAJOR_VERSION.x.x automatically)"
echo "  @v$NEW_VERSION (pinned to this exact version)"
