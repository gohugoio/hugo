---
title: default version
---

The _default version_ is the value defined by the [`defaultContentVersion`][] setting, falling back to `v1.0.0` when the project defines no versions. When the project defines one or more versions and the setting is not defined, the default version is the first version in the project, identified by the lowest [_weight_](g), using a descending semantic sort as the final fallback if weights are tied or undefined.

  See also: [_version_](g).

  [`defaultContentVersion`]: /configuration/all/#defaultcontentversion
