#!/bin/sh

set -ex

export DART_SASS_VERSION=1.79.3

# If $BUILDARCH=arm64, then we need to install the arm64 version of Dart Sass,
# otherwise we install the x64 version.
ARCH="x64"
if [ "$BUILDARCH" = "arm64" ]; then
    ARCH="arm64"
fi

cd /tmp
curl -LJO https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-${ARCH}.tar.gz 
ls -ltr
tar -xf dart-sass-${DART_SASS_VERSION}-linux-${ARCH}.tar.gz
rm dart-sass-${DART_SASS_VERSION}-linux-${ARCH}.tar.gz && \
# The dart-sass folder is added to the PATH by the caller.
mv dart-sass /var/hugo/bin