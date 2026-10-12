#!/bin/bash

set -e

# The Javy version used to build the QuickJS plugin and the JS WASM files.
# Download releases from https://github.com/bytecodealliance/javy/releases.
JAVY_VERSION="9.1.0"

v=$(javy -V 2>/dev/null | awk '{print $2}')
if [ "$v" != "$JAVY_VERSION" ]; then
	echo "error: javy $JAVY_VERSION required, found '$v'"
	exit 1
fi

go generate ./genjs
javy emit-plugin -o wasm/quickjs.wasm
javy build -C dynamic -C plugin=wasm/quickjs.wasm -o wasm/greet.wasm js/greet.bundle.js
javy build -C dynamic -C plugin=wasm/quickjs.wasm -o wasm/renderkatex.wasm js/renderkatex.bundle.js
touch warpc_test.go
