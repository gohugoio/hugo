# GitHub:       https://github.com/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.5.0 AS xx

FROM --platform=$BUILDPLATFORM golang:1.22.6-alpine AS build

# Set up cross-compilation helpers
COPY --from=xx / /
RUN apk add clang lld

# Optionally set HUGO_BUILD_TAGS to "extended" or "nodeploy" when building like so:
#   docker build --build-arg HUGO_BUILD_TAGS=extended .
ARG HUGO_BUILD_TAGS="none"

ARG CGO=1
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/gohugoio/hugo

RUN --mount=src=go.mod,target=go.mod \
    --mount=src=go.sum,target=go.sum \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

ARG TARGETPLATFORM
# gcc/g++ are required to build SASS libraries for extended version
RUN xx-apk add --no-scripts --no-cache gcc g++ musl-dev git
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod <<EOT
    set -ex
    xx-go build -tags "$HUGO_BUILD_TAGS" -o /usr/bin/hugo
    xx-verify /usr/bin/hugo
EOT

# ---

FROM alpine:3.18

COPY --from=build /usr/bin/hugo /usr/bin/hugo

# libc6-compat & libstdc++ are required for extended SASS libraries
# ca-certificates are required to fetch outside resources (like Twitter oEmbeds)
RUN apk update && \
    apk add --no-cache ca-certificates libc6-compat libstdc++ git

VOLUME /site
WORKDIR /site

# Expose port for live server
EXPOSE 1313

ENTRYPOINT ["hugo"]
CMD ["--help"]
