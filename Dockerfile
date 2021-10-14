# GitHub:       https://github.com/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

# crane digest golang:1.16-alpine https://github.com/google/go-containerregistry/blob/main/cmd/crane/README.md
# sha256:78181bcf43be59a818e23095f21b3818f456895c3f1f2daaabdfd1af75cedd1f
# Pinning by SHA
# The original version before the SHA was golang:1.16-alpine

FROM golang@sha256:78181bcf43be59a818e23095f21b3818f456895c3f1f2daaabdfd1af75cedd1f AS build

# Optionally set HUGO_BUILD_TAGS to "extended" or "nodeploy" when building like so:
#   docker build --build-arg HUGO_BUILD_TAGS=extended .
ARG HUGO_BUILD_TAGS

ARG CGO=1
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/gohugoio/hugo

COPY . /go/src/github.com/gohugoio/hugo/

# gcc/g++ are required to build SASS libraries for extended version
RUN apk update && \
    apk add --no-cache gcc g++ musl-dev && \
    go get github.com/magefile/mage

RUN mage hugo && mage install

# ---

# crane digest alpine:3.12
# sha256:a296b4c6f6ee2b88f095b61e95c7ef4f51ba25598835b4978c9256d8c8ace48a
# Pinning by SHA
# The original version before the SHA was alpine:3.12

FROM alpine@sha256:a296b4c6f6ee2b88f095b61e95c7ef4f51ba25598835b4978c9256d8c8ace48a

COPY --from=build /go/bin/hugo /usr/bin/hugo

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
