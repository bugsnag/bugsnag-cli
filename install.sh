#!/usr/bin/env bash

set -x

# Detect the version of OS we're using
UNAME=$(uname)
# Detect the arch type that is being used
ARCH=$(uname -m)

if [ "$UNAME" == "Linux" ] ; then
  OS="linux"
elif [ "$UNAME" == "Darwin" ] ; then
  OS="macos"
fi

echo "Downloading $ARCH-$OS-bugsnag-cli"

# Get download URL
DOWNLOAD_URL=$(curl -s https://api.github.com/repos/joshedney/red-bucket/releases/latest | \
grep "$ARCH-$OS-bugsnag-cli*" | \
grep "browser_download_url" | \
cut -d : -f 2,3 | \
tr -d \")

# Download release
curl -L $DOWNLOAD_URL -o /tmp/bugsnag-cli

if [[ "$OS" == "linux" || "$OS" == "macos" ]]; then
  echo "Asking for sudo access"
  sudo mv /tmp/bugsnag-cli /usr/local/bin/bugsnag-cli
  sudo chmod +x /usr/local/bin/bugsnag-cli
  echo "bugsnag-cli added to /usr/local/bin"
fi

# Remove app from quarantine on macos
if [ "$OS" == "macos" ]; then
  echo "Removing file from Apple quarantine"
  xattr -d com.apple.quarantine /usr/local/bin/bugsnag-cli
fi
