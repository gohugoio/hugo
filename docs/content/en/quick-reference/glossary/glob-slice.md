---
title: glob slice
---

A _glob slice_ is a [_slice_](g) of [_glob patterns_](g). Within the _slice_, a _glob_ can be negated by prefixing it with an exclamation mark (`!`) and one space. Matches in negated patterns short-circuit the evaluation of the rest of the _slice_, and are useful for early coarse grained exclusions.

  The following example illustrates how to use _glob slices_ to define a [_sites matrix_](g) in your project configuration:

  ```toml
  [sites.matrix]
  languages = [ "! no", "**" ]
  versions = [ "! v1.2.3", "v1.*.*", "v2.*.*" ]
  roles = [ "{member, guest}" ]
  ```

  The `versions` example above evaluates as: `(not v1.2.3) AND (v1.*.* OR v2.*.*)`.
