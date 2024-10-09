# GitHub:       https://github.com/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.5.0 AS xx
FROM --platform=$BUILDPLATFORM golang:1.23.2-alpine AS base
FROM base AS build

# Set up cross-compilation helpers
COPY --from=xx / /

# gcc/g++ are required to build libsass and libwebp libraries for the extended version.
RUN xx-apk add --no-scripts --no-cache gcc g++

# Optionally set HUGO_BUILD_TAGS to "none" or "nodeploy" when building like so:
# docker build --build-arg HUGO_BUILD_TAGS=nodeploy .
#
# We build the extended version by default.
ARG HUGO_BUILD_TAGS="extended"
ENV GOPROXY=https://proxy.golang.org
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/go/pkg/mod

WORKDIR /go/src/github.com/gohugoio/hugo

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build <<EOT
    set -ex
    xx-go build -tags "$HUGO_BUILD_TAGS" -ldflags "-s -w -X github.com/gohugoio/hugo/common/hugo.vendorInfo=docker" -o /usr/bin/hugo
    xx-verify /usr/bin/hugo
EOT

FROM base AS final

COPY --from=build /usr/bin/hugo /usr/bin/hugo

# libc6-compat & libstdc++ are required for extended libraries (libsass, libwebp).
RUN apk add --no-cache \
    libc6-compat \
    libstdc++ \
    git \
    runuser \
    curl \
    nodejs \
    npm

RUN mkdir -p /var/hugo/bin && \
    addgroup -Sg 1000 hugo && \
    adduser -Sg hugo -u 1000 -h /var/hugo hugo && \
    chown -R hugo: /var/hugo && \
    runuser -u hugo -- git config --global --add safe.directory /project

VOLUME /project
WORKDIR /project
USER hugo:hugo
ENV HUGO_CACHEDIR=/cache
ARG BUILDARCH
ENV BUILDARCH=${BUILDARCH}

COPY scripts/docker scripts/docker

# Install default dependencies.
RUN scripts/docker/install_runtimedeps_default.sh

COPY scripts/docker/entrypoint.sh /entrypoint.sh

ENV PATH="/var/hugo/bin:/var/hugo/bin/dart-sass:$PATH"

RUN sass --version
RUN hugo version

# Expose port for live server
EXPOSE 1313

ENTRYPOINT ["/entrypoint.sh"]
CMD ["--help"]

