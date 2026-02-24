---
title: default role
---

The _default role_ is the value defined by the [`defaultContentRole`][] setting, falling back to the first role in the project, and finally to `guest`. The first role is identified by the lowest [_weight_](g), using lexicographical order as the final fallback if weights are tied or undefined.

  See also: [_role_](g).

  [`defaultContentRole`]: /configuration/all/#defaultcontentrole
