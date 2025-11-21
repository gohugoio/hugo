# GitHub:       https://github.com/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

ARG GO_VERSION="1.25"
ARG ALPINE_VERSION="3.22"
ARG DART_SASS_VERSION="1.79.3"

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.5.0 AS xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS gobuild
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS gorun

FROM gobuild AS build

# Required tools for CGO + cross-compilation
RUN apk add --no-cache clang lld

# Set up cross-compilation helpers
COPY --from=xx / /

ARG TARGETPLATFORM
RUN xx-apk add --no-cache musl-dev gcc g++

# Build tags (extended version by default)
ARG HUGO_BUILD_TAGS="extended"

# Consolidated ENV block
ENV \
  CGO_ENABLED=1 \
  GOPROXY=https://proxy.golang.org \
  GOCACHE=/root/.cache/go-build \
  GOMODCACHE=/go/pkg/mod

WORKDIR /go/src/github.com/gohugoio/hugo

# Build Hugo with caching
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build,id=go-build-$TARGETPLATFORM <<EOT
    set -ex
    xx-go build -tags "$HUGO_BUILD_TAGS" \
      -ldflags "-s -w -X github.com/gohugoio/hugo/common/hugo.vendorInfo=docker" \
      -o /usr/bin/hugo
    xx-verify /usr/bin/hugo
EOT

# Download + extract Dart Sass
FROM alpine:${ALPINE_VERSION} AS dart-sass
ARG TARGETARCH
ARG DART_SASS_VERSION
ARG DART_ARCH=${TARGETARCH/amd64/x64}

WORKDIR /out

# Use wget instead of ADD, then cleanup
RUN apk add --no-cache wget && \
    wget -q https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-${DART_ARCH}.tar.gz && \
    tar -xf dart-sass-${DART_SASS_VERSION}-linux-${DART_ARCH}.tar.gz && \
    rm dart-sass-${DART_SASS_VERSION}-linux-${DART_ARCH}.tar.gz

FROM gorun AS final

COPY --from=build /usr/bin/hugo /usr/bin/hugo

# Required runtime dependencies
RUN apk add --no-cache \
    libc6-compat \
    git \
    runuser \
    nodejs \
    npm

# Create user, directories, and permissions
RUN mkdir -p /var/hugo/bin /cache && \
    addgroup -g 1000 -S hugo && \
    adduser -u 1000 -G hugo -S -h /var/hugo hugo && \
    chown -R hugo: /var/hugo /cache && \
    runuser -u hugo -- git config --global --add safe.directory /project && \
    runuser -u hugo -- git config --global core.quotepath false

USER hugo:hugo
VOLUME /project
WORKDIR /project

# Consolidated ENV block
ENV \
  HUGO_CACHEDIR=/cache \
  PATH="/var/hugo/bin:/var/hugo/bin/dart-sass:$PATH"

COPY scripts/docker/entrypoint.sh /entrypoint.sh
COPY --from=dart-sass /out/dart-sass /var/hugo/bin/dart-sass

# Expose port for live server
EXPOSE 1313

ENTRYPOINT ["/entrypoint.sh"]
CMD ["--help"]
