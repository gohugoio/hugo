# The library versions to use.
LIBAVIF_VERSION = v1.4.1
AOM_VERSION = v3.14.1
DAV1D_VERSION = 1.5.3
# Commit used by Chromium, see https://source.chromium.org/chromium/chromium/src/+/main:DEPS?q=libyuv
LIBYUV_COMMIT = 9d98aaefe7a5e2710aa6175d44da38892400b381

# The path to the WebAssembly SDK. Can be overridden from the command line.
# e.g., make WASI_SDK_PATH=/path/to/wasi-sdk
# Download releases from https://github.com/WebAssembly/wasi-sdk/releases.
# Note that on MacOS, you may need to add some whitelistings in Privacy & Security > Security.
WASI_SDK_PATH ?= /opt/wasi-sdk

LIBAVIF_BUILD := $(PWD)/build
LIBAVIF_SRC := $(PWD)/libavif
LIBAVIF_AOM_SRC := $(LIBAVIF_SRC)/ext/aom
LIBAVIF_AOM_BUILD = $(LIBAVIF_AOM_SRC)/build.libavif
LIBAVIF_DAV1D_SRC := $(LIBAVIF_SRC)/ext/dav1d
LIBAVIF_DAV1D_BUILD := $(LIBAVIF_DAV1D_SRC)/build
LIBYUV_SRC := $(LIBAVIF_SRC)/ext/libyuv
LIBYUV_BUILD := $(LIBYUV_SRC)/build

BUILD_WASM := ${LIBAVIF_BUILD}/avif.wasm
TARGET_WASM_DIR := ../wasm

# Set the compiler and toolchain for WASI
CC := $(WASI_SDK_PATH)/bin/clang
CMAKE_TOOLCHAIN_FILE := $(WASI_SDK_PATH)/share/cmake/wasi-sdk-p1.cmake

# Source files for the final WASM binary
C_SOURCES := avif.c ../deps/parson/parson.c

# Include paths
C_INCLUDES := -I$(LIBAVIF_SRC)/include -I$(LIBAVIF_AOM_SRC) -I$(LIBAVIF_BUILD)/include -I.

# Libraries to link against
LD_LIBS := $(LIBAVIF_BUILD)/libavif.a $(LIBAVIF_DAV1D_BUILD)/src/libdav1d.a $(LIBAVIF_BUILD)/libaom.a $(LIBYUV_BUILD)/libyuv.a -lsetjmp

# SIMD + non-trapping FP/int are part of wazero's CoreFeaturesV2, so we can rely
# on them at runtime. Adding -msimd128 to the library builds lets clang
# auto-vectorize hot loops in aom/dav1d/libyuv/libavif (none of them ship
# hand-tuned wasm-simd code).
LIB_SIMD_FLAGS := -msimd128 -mnontrapping-fptoint
LIB_CFLAGS := $(LIB_SIMD_FLAGS) -mllvm -wasm-enable-sjlj -mllvm -wasm-use-legacy-eh=false

# Compiler flags for the WASM build
WASM_CFLAGS := \
    -O3 \
    -msimd128 \
    -mexec-model=command \
    -mnontrapping-fptoint \
    -mllvm -wasm-enable-sjlj \
    -mllvm -wasm-use-legacy-eh=false \
    -Wl,--export=malloc \
    -Wl,--export=free \
    -Wall \

all: $(LIBAVIF_SRC) $(LIBAVIF_AOM_SRC) $(LIBAVIF_DAV1D_SRC) $(LIBYUV_SRC) $(LIBAVIF_BUILD)
	@echo ">>> Configuring and building C libraries in $(LIBAVIF_BUILD)"
	cd $(LIBAVIF_AOM_BUILD); \
	cmake $(LIBAVIF_AOM_SRC) \
		-DCMAKE_BUILD_TYPE=Release \
		-DBUILD_SHARED_LIBS=0 \
		-DENABLE_DOCS=0 \
		-DENABLE_EXAMPLES=0 \
		-DENABLE_TESTDATA=0 \
		-DENABLE_TESTS=0 \
		-DENABLE_TOOLS=0 \
		-DAOM_TARGET_CPU=generic \
		-DCONFIG_RUNTIME_CPU_DETECT=0 \
		-DCONFIG_MULTITHREAD=0 \
		-DCONFIG_WEBM_IO=0 \
		-DCONFIG_AV1_DECODER=1 \
		-DCONFIG_AV1_ENCODER=1 \
		-DCMAKE_TOOLCHAIN_FILE=$(CMAKE_TOOLCHAIN_FILE) \
		-DCMAKE_C_FLAGS="$(LIB_CFLAGS)" \
		-DCMAKE_CXX_FLAGS="$(LIB_CFLAGS)"
	cd $(LIBAVIF_AOM_BUILD); \
	make VERBOSE=1
	@echo ">>> Copying libaom.a to $(LIBAVIF_BUILD)"
	cp $(LIBAVIF_AOM_BUILD)/libaom.a $(LIBAVIF_BUILD)/libaom.a

	@echo ">>> Building dav1d with meson"
	cd $(LIBAVIF_DAV1D_SRC); \
	meson setup build \
		--cross-file $(PWD)/wasi-sdk-cross.txt \
		-Denable_asm=false \
		-Denable_tools=false \
		-Denable_tests=false \
		-Ddefault_library=static \
		--reconfigure
	cd $(LIBAVIF_DAV1D_SRC); \
	ninja -C build -v

	@echo ">>> Building libyuv"
	mkdir -p $(LIBYUV_BUILD)
	cd $(LIBYUV_BUILD); \
	cmake $(LIBYUV_SRC) \
		-DCMAKE_BUILD_TYPE=Release \
		-DBUILD_SHARED_LIBS=0 \
		-DCMAKE_POSITION_INDEPENDENT_CODE=ON \
		-DCMAKE_TOOLCHAIN_FILE=$(CMAKE_TOOLCHAIN_FILE) \
		-DCMAKE_POLICY_VERSION_MINIMUM=3.5 \
		-DCMAKE_C_FLAGS="$(LIB_CFLAGS)" \
		-DCMAKE_CXX_FLAGS="$(LIB_CFLAGS)"
	cd $(LIBYUV_BUILD); \
	make yuv VERBOSE=1

	cd $(LIBAVIF_BUILD); \
	cmake $(LIBAVIF_SRC) \
		-DCMAKE_BUILD_TYPE=Release \
		-DBUILD_SHARED_LIBS=0 \
		-DAVIF_CODEC_AOM=LOCAL \
		-DAVIF_LOCAL_AOM=1 \
		-DAVIF_CODEC_AOM_DECODE=1 \
		-DAVIF_CODEC_AOM_ENCODE=1 \
		-DAVIF_CODEC_DAV1D=LOCAL \
		-DAVIF_LIBYUV=LOCAL \
		-DCMAKE_TOOLCHAIN_FILE=$(CMAKE_TOOLCHAIN_FILE) \
		-DCMAKE_C_FLAGS="$(LIB_CFLAGS)" \
		-DCMAKE_CXX_FLAGS="$(LIB_CFLAGS)"

	cd $(LIBAVIF_BUILD); \
	make VERBOSE=1

	@echo ">>> Building WebAssembly binary: $(BUILD_WASM)"
	$(CC) $(WASM_CFLAGS) --sysroot=$(WASI_SDK_PATH)/share/wasi-sysroot $(C_INCLUDES) -o $(BUILD_WASM) $(C_SOURCES) $(LD_LIBS)
	@echo ">>> Moving the WebAssembly binary: $(TARGET_WASM_DIR)"
	mv $(BUILD_WASM) $(TARGET_WASM_DIR)


# Rule to create the build directory
$(LIBAVIF_BUILD):
	@echo ">>> Creating build directory: $(LIBAVIF_BUILD)"
	@mkdir -p $@

# Rule to clone the libavif source code
$(LIBAVIF_SRC):
	@echo ">>> Cloning libavif source code (version: $(LIBAVIF_VERSION))"
	git clone -b $(LIBAVIF_VERSION) --depth 1 --recursive https://github.com/AOMediaCodec/libavif.git  $(LIBAVIF_SRC)
	test -d $@

# Rule to clone the aom source code
$(LIBAVIF_AOM_SRC):
	@echo ">>> Cloning aom source code (version: $(AOM_VERSION))"
	git clone -b $(AOM_VERSION) --depth 1 --recursive https://aomedia.googlesource.com/aom $(LIBAVIF_AOM_SRC)
	mkdir -p $(LIBAVIF_AOM_BUILD)
	test -d $@

# Rule to clone the dav1d source code
$(LIBAVIF_DAV1D_SRC):
	@echo ">>> Cloning dav1d source code (version: $(DAV1D_VERSION))"
	git clone -b $(DAV1D_VERSION) --depth 1 https://code.videolan.org/videolan/dav1d.git $(LIBAVIF_DAV1D_SRC)
	test -d $@

# Rule to clone the libyuv source code
$(LIBYUV_SRC):
	@echo ">>> Cloning libyuv source code (commit: $(LIBYUV_COMMIT))"
	git clone --single-branch https://chromium.googlesource.com/libyuv/libyuv $(LIBYUV_SRC)
	cd $(LIBYUV_SRC) && git checkout $(LIBYUV_COMMIT)
	test -d $@


.PHONY: clean

# Rule to clean up build artifacts
clean:
	@echo ">>> Cleaning build artifacts"
	@rm -rf $(LIBAVIF_BUILD)
	@rm -rf $(LIBAVIF_SRC)
