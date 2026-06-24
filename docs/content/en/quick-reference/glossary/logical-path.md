---
title: logical path
reference: /methods/page/path/#logical-tree
---

A _logical path_ is a page or page resource identifier derived from the file path, excluding its extension and language identifier. This value is neither a file path nor a URL. Starting with a file path relative to the `content` directory, Hugo determines the logical path by stripping the file extension and language identifier, converting to lower case, then replacing spaces with hyphens. Path segments are separated with a slash (`/`).

  When used to describe content, the logical path is the path between two [_nodes_](g) in the [_logical tree_](g), either relative to each other or absolute from the root of the tree.
