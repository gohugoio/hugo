---
title: site root
---

The _site root_ is the root directory of the current [_site_](g), relative to the [`publishDir`][]. The _site root_ may include one or more content [_dimension_](g) prefixes, such as [_language_](g), [_role_](g), or [_version_](g).

  Project description|Site root examples
  :--|:--|:--
  Monolingual|`/`, `/guest`, `/guest/v1.2.3`
  Multilingual single-host|`/en`, `/guest/en`, `/guest/v1.2.3/en`
  Multilingual multihost|`/en`, `/en/guest`, `/en/guest/v1.2.3`

  [`publishDir`]: /configuration/all/#publishdir
