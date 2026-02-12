#!/bin/bash
# bump-version.sh: Increment SemVer version strings

set -e

if [ $# -lt 2 ]; then
    echo "Usage: $0 <current_version> <bump_type>"
    echo "bump_type: patch, minor, major"
    exit 1
fi

CURRENT_VERSION=$1
BUMP_TYPE=$2

# Strip leading 'v' if present
VERSION=${CURRENT_VERSION#v}

# Split version into components
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

case "$BUMP_TYPE" in
    patch)
        PATCH=$((PATCH + 1))
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    *)
        echo "Error: Unknown bump type '$BUMP_TYPE'"
        exit 1
        ;;
esac

echo "v$MAJOR.$MINOR.$PATCH"
