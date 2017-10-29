#!/bin/bash

HUGO_DOCS_BRANCH="${HUGO_DOCS_BRANCH-master}"

# We may extend this to also push changes in the other direction, but this is the most important step.
git subtree pull --prefix=docs/ https://github.com/gohugoio/hugoDocs.git ${HUGO_DOCS_BRANCH} --squash

