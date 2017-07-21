#!/bin/bash

# We may extend this to also push changes in the other direction, but this is the most important step.
git subtree pull --prefix=docs/ https://github.com/gohugoio/hugoDocs.git master --squash

