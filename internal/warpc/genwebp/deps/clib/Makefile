CC     ?= cc
PREFIX ?= /usr/local

BINS = clib clib-install clib-search clib-init clib-configure clib-build clib-update clib-upgrade clib-uninstall

ifdef EXE
	BINS := $(addsuffix .exe,$(BINS))
endif

CP      = cp -f
RM      = rm -f
MKDIR   = mkdir -p
INSTALL = install

SRC  = $(wildcard src/*.c)
COMMON_SRC = $(wildcard src/common/*.c)
ALL_SRC = $(wildcard src/*.c src/*.h src/common/*.c src/common/*.h test/package/*.c test/cache/*.c)
SDEPS = $(wildcard deps/*/*.c)
ODEPS = $(SDEPS:.c=.o)
DEPS = $(filter-out $(ODEPS), $(SDEPS))
OBJS = $(DEPS:.c=.o)
MAKEFILES = $(wildcard deps/*/Makefile)
HEADERS_BINS = src/common/*.h src/version.h deps/logger/logger.h

export CC

CFLAGS += -std=c99 -Ideps -Wall -Wno-unused-function -U__STRICT_ANSI__

ifdef STATIC
	CFLAGS  += -DCURL_STATICLIB  $(shell deps/curl/bin/curl-config --cflags)
	LDFLAGS += -static $(shell deps/curl/bin/curl-config --static-libs)
else
	CFLAGS  += $(shell curl-config --cflags)
	LDFLAGS += $(shell curl-config --libs)
endif

ifneq (0,$(PTHREADS))
ifndef NO_PTHREADS
	CFLAGS += $(shell ./scripts/feature-test-pthreads && echo "-DHAVE_PTHREADS=1 -pthread")
endif
endif

ifdef DEBUG
	CFLAGS += -g -D CLIB_DEBUG=1 -D DEBUG="$(DEBUG)"
endif

default: all

all: $(BINS)

build: $(BINS)

$(BINS): $(SRC) $(COMMON_SRC) $(MAKEFILES) $(OBJS) $(HEADERS_BINS)
	$(CC) $(CFLAGS) -o $@ $(COMMON_SRC) src/$(@:.exe=).c $(OBJS) $(LDFLAGS)

$(MAKEFILES):
	$(MAKE) -C $@

%.o: %.c
	$(CC) $< -c -o $@ $(CFLAGS) -MMD

clean:
	$(foreach c, $(BINS), $(RM) $(c);)
	$(RM) $(OBJS)
	$(RM) $(AUTODEPS)
	cd test/cache && make clean
	cd test/package && make clean

install: $(BINS)
	$(MKDIR) $(PREFIX)/bin
	$(foreach c, $(BINS), $(INSTALL) $(c) $(PREFIX)/bin/$(c);)

uninstall:
	$(foreach c, $(BINS), $(RM) $(PREFIX)/bin/$(c);)

test: $(BINS)
	@./test.sh

# create a list of auto dependencies
AUTODEPS:= $(patsubst %.c,%.d, $(DEPS)) $(patsubst %.c,%.d, $(SRC))

# include by auto dependencies
-include $(AUTODEPS)

# Format all source files in the repository.
fmt:
	@if ! command -v clang-format &> /dev/null; then \
		echo "clang-format not found"; \
		exit; \
	fi
	clang-format -i -style=LLVM $(ALL_SRC)

# Install the commit hook.
commit-hook: scripts/pre-commit-hook.sh
	cp -f scripts/pre-commit-hook.sh .git/hooks/pre-commit

.PHONY: test all clean install uninstall fmt
