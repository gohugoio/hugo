docker run --rm --mount type=bind,source="$(pwd)",target=hugo -w /hugo  -i -t bepsays/ci-goreleaser:1.11-2 /bin/bash
