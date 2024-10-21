# GitHub:       https://github.com/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

ARG GO_VERSION="1.23.2"
ARG ALPINE_VERSION="3.20"
ARG DART_SASS_VERSION="1.79.3"

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.5.0 AS xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS gobuild
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS gorun


FROM gobuild AS build

RUN apk add clang lld

# Set up cross-compilation helpers
COPY --from=xx / /

ARG TARGETPLATFORM
RUN xx-apk add musl-dev gcc g++ 

# Optionally set HUGO_BUILD_TAGS to "none" or "nodeploy" when building like so:
# docker build --build-arg HUGO_BUILD_TAGS=nodeploy .
#
# We build the extended version by default.
ARG HUGO_BUILD_TAGS="extended"
ENV CGO_ENABLED=1
ENV GOPROXY=https://proxy.golang.org
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/go/pkg/mod
ARG TARGETPLATFORM

WORKDIR /go/src/github.com/gohugoio/hugo

# For  --mount=type=cache the value of target is the default cache id, so
# for the go mod cache it would be good if we could share it with other Go images using the same setup,
# but the go build cache needs to be per platform.
# See this comment: https://github.com/moby/buildkit/issues/1706#issuecomment-702238282
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build,id=go-build-$TARGETPLATFORM <<EOT
    set -ex
    xx-go build -tags "$HUGO_BUILD_TAGS" -ldflags "-s -w -X github.com/gohugoio/hugo/common/hugo.vendorInfo=docker" -o /usr/bin/hugo
    xx-verify /usr/bin/hugo
EOT

# dart-sass downloads the dart-sass runtime dependency
FROM alpine:${ALPINE_VERSION} AS dart-sass
ARG TARGETARCH
ARG DART_SASS_VERSION
ARG DART_ARCH=${TARGETARCH/amd64/x64}
WORKDIR /out
ADD https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-${DART_ARCH}.tar.gz .
RUN tar -xf dart-sass-${DART_SASS_VERSION}-linux-${DART_ARCH}.tar.gz

FROM gorun AS final

COPY --from=build /usr/bin/hugo /usr/bin/hugo

# libc6-compat  are required for extended libraries (libsass, libwebp).
RUN apk add --no-cache \
    libc6-compat \
    git \
    runuser \
    nodejs \
    npm

RUN mkdir -p /var/hugo/bin /cache && \
    addgroup -Sg 1000 hugo && \
    adduser -Sg hugo -u 1000 -h /var/hugo hugo && \
    chown -R hugo: /var/hugo /cache && \
    # For the Hugo's Git integration to work.
    runuser -u hugo -- git config --global --add safe.directory /project && \ 
    # See https://github.com/gohugoio/hugo/issues/9810
    runuser -u hugo -- git config --global core.quotepath false

USER hugo:hugo
VOLUME /project
WORKDIR /project
ENV HUGO_CACHEDIR=/cache
ENV PATH="/var/hugo/bin:$PATH"

COPY scripts/docker/entrypoint.sh /entrypoint.sh
COPY --from=dart-sass /out/dart-sass /var/hugo/bin/dart-sass

# Update PATH to reflect the new dependencies.
# For more complex setups, we should probably find a way to
# delegate this to the script itself, but this will have to do for now.
# Also, the dart-sass binary is a little special, other binaries can be put/linked
# directly in /var/hugo/bin.
ENV PATH="/var/hugo/bin/dart-sass:$PATH"

# Expose port for live server
EXPOSE 1313

ENTRYPOINT ["/entrypoint.sh"]
CMD ["--help"]
