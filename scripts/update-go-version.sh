#!/bin/bash
# Script to update Go version across all project files

if [ $# -ne 1 ]; then
    echo "Usage: $0 <go-version>"
    echo "Example: $0 1.21.0"
    exit 1
fi

NEW_VERSION=$1
MAJOR_MINOR=$(echo $NEW_VERSION | cut -d. -f1-2)

echo "Updating Go version to $NEW_VERSION..."

# Update .go-version
echo "$NEW_VERSION" > .go-version

# Update go.mod
sed -i '' "s/^go .*/go $MAJOR_MINOR/" go.mod

# Update Dockerfile
sed -i '' "s/golang:[0-9.]*-alpine/golang:$MAJOR_MINOR-alpine/" Dockerfile

echo "Updated Go version to $NEW_VERSION in:"
echo "  - .go-version"
echo "  - go.mod (using $MAJOR_MINOR)"
echo "  - Dockerfile (using $MAJOR_MINOR-alpine)"
echo ""
echo "GitHub Actions will automatically use the version from .go-version"
