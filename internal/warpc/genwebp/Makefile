# The library version to use.
# Note that when updating this, you may also need to update the badge in README.md.
LIBWEBP_VERSION = v1.6.0

# The path to the WebAssembly SDK. Can be overridden from the command line.
# e.g., make WASI_SDK_PATH=/path/to/wasi-sdk
# Download releases from https://github.com/WebAssembly/wasi-sdk/releases.
# Note that on MacOS, you may need to add some whitelistings in Privacy & Security > Security.
WASI_SDK_PATH ?= /opt/wasi-sdk

BUILD_DIR := build
BUILD_WASM := ${BUILD_DIR}/webp.wasm
TARGET_WASM_DIR := ../wasm
LIBWEBP_SRC := libwebp

# Set the compiler and toolchain for WASI
CC := $(WASI_SDK_PATH)/bin/clang
CMAKE_TOOLCHAIN_FILE := $(WASI_SDK_PATH)/share/cmake/wasi-sdk.cmake

# Source files for the final WASM binary
C_SOURCES := webp.c deps/parson/parson.c

# Include paths
C_INCLUDES := -I$(LIBWEBP_SRC)/src -I$(BUILD_DIR)/src -I.

# Libraries to link against
LD_LIBS := $(BUILD_DIR)/libwebpdemux.a $(BUILD_DIR)/libwebp.a $(BUILD_DIR)/libsharpyuv.a $(BUILD_DIR)/libwebpmux.a

# Compiler flags for the WASM build
WASM_CFLAGS := \
    -O3 \
    -msimd128 \
    -mexec-model=command \
    -mnontrapping-fptoint \
    -Wl,--export=malloc \
    -Wl,--export=free \
    -Wall

.PHONY: all clean

# Default target: build the WebAssembly binary
all: $(LIBWEBP_SRC) $(BUILD_DIR)
	@echo ">>> Configuring and building C libraries in $(BUILD_DIR)"
	cd $(BUILD_DIR) && \
	cmake $(realpath $(LIBWEBP_SRC)) \
		-DCMAKE_TOOLCHAIN_FILE=$(CMAKE_TOOLCHAIN_FILE) \
		-DCMAKE_BUILD_TYPE=Release \
		-DBUILD_SHARED_LIBS=0 \
		-DWEBP_ENABLE_SIMD_DEFAULT=1 \
		-DWEBP_BUILD_EXTRAS=0 \
		-DWEBP_USE_THREAD=0 \
		-DWEBP_BUILD_ANIM_UTILS=0 \
		-DWEBP_BUILD_CWEBP=0 \
		-DWEBP_BUILD_DWEBP=0 \
		-DWEBP_BUILD_IMG2WEBP=0 \
		-DWEBP_BUILD_WEBPINFO=0 \
		-DWEBP_BUILD_WEBPMUX=1 && \
	make VERBOSE=1
	@echo ">>> Building WebAssembly binary: $(BUILD_WASM)"
	$(CC) $(WASM_CFLAGS) --sysroot=$(WASI_SDK_PATH)/share/wasi-sysroot $(C_INCLUDES) -o $(BUILD_WASM) $(C_SOURCES) $(LD_LIBS)
	@echo ">>> Moving the WebAssembly binary: $(TARGET_WASM_DIR)"
	mv $(BUILD_WASM) $(TARGET_WASM_DIR)



# Rule to create the build directory
$(BUILD_DIR):
	@echo ">>> Creating build directory: $(BUILD_DIR)"
	@mkdir -p $@

# Rule to clone the libwebp source code
$(LIBWEBP_SRC):
	@echo ">>> Cloning libwebp source code (version: $(LIBWEBP_VERSION))"
	git clone -b $(LIBWEBP_VERSION) --depth 1 --recursive https://github.com/webmproject/libwebp
	test -d $@

# Rule to clean up build artifacts
clean:
	@echo ">>> Cleaning build artifacts"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(LIBWEBP_SRC)
