---
title: default role
---

The _default role_ is the value defined by the [`defaultContentRole`][] setting, falling back to `guest` when the project defines no roles. When the project defines one or more roles and the setting is not defined, the default role is the first role in the project, identified by the lowest [_weight_](g), using lexicographical order as the final fallback if weights are tied or undefined.

  See also: [_role_](g).

  [`defaultContentRole`]: /configuration/all/#defaultcontentrole
