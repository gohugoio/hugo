---
title: default version
---

The _default version_ is the value defined by the [`defaultContentVersion`][] setting, falling back to the first version in the project, and finally to `v1.0.0`. The first version is identified by the lowest [_weight_](g), using a descending semantic sort as the final fallback if weights are tied or undefined.

  See also: [_version_](g).

  [`defaultContentVersion`]: /configuration/all/#defaultcontentversion
