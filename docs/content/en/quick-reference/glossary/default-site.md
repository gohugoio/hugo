---
title: default site
---

The _default site_ is the [_site_](g) identified by the primary value in each [_content dimension_](g). Specifically, it is the site that combines the first language, the first [_role_](g), and the first version defined in your site configuration.

  The "first" language and role are those with the lowest [weight](g). If weights are tied or undefined, Hugo defaults to lexicographical order. Similarly, the "first" version is the one with the lowest weight; if weights are tied or undefined, it is identified as the last version when sorted semantically.
