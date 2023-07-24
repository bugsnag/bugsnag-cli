#!/bin/bash -e

BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$1" != "" ]]; then
  VERSION=$1
elif [[ "$BRANCH" =~ ^release\/v.*$ ]]; then
  VERSION=${BRANCH#*release/v}
else
  echo "$BRANCH does not appear to be a release branch, please specify VERSION manually"
  echo "$(basename $0) <version-number>"
  exit 1
fi

echo Bumping the version number to $VERSION
sed -i '' "s/package_version = \".*\"/package_version = \"$VERSION\"/" main.go
sed -i '' "s/\"version\": \".*\"/\"version\": \"$VERSION\"/" package.json
sed -i '' "s/# TBD/# $VERSION ($(date '+%Y-%m-%d'))/" CHANGELOG.md
