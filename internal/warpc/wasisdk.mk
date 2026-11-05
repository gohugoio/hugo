# The wasi-sdk major version used to build the WASM binaries.
# Download releases from https://github.com/WebAssembly/wasi-sdk/releases.
WASI_SDK_VERSION = 34

# The path to the WebAssembly SDK. Can be overridden from the command line.
# e.g., make WASI_SDK_PATH=/path/to/wasi-sdk
# Note that on MacOS, you may need to add some whitelistings in Privacy & Security > Security.
WASI_SDK_PATH ?= /opt/wasi-sdk

# Flags for linking the final WASM binaries.
#
# clang runs wasm-opt on optimized builds if it finds it in PATH (e.g. from a
# Homebrew binaryen install), which strips the name section, so function
# names would be missing from the stack traces of WASM errors, and the output
# would depend on what's installed on the build machine. wasm-opt doesn't make
# a measurable difference at runtime, as wazero optimizes the code when
# compiling it.
WASM_LDFLAGS_COMMON := --no-wasm-opt

# Keep "all" as the default goal in the including Makefiles.
.DEFAULT_GOAL := all

.PHONY: check-wasi-sdk
check-wasi-sdk:
	@v=$$(head -n 1 "$(WASI_SDK_PATH)/VERSION" 2>/dev/null); \
	if [ "$${v%%.*}" != "$(WASI_SDK_VERSION)" ]; then \
		echo "error: wasi-sdk $(WASI_SDK_VERSION) required, found '$$v' in $(WASI_SDK_PATH)"; \
		exit 1; \
	fi
