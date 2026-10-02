---
title: default language
---

The _default language_ is the value defined by the [`defaultContentLanguage`][] setting, falling back to `en` when the project defines no languages. When the project defines one or more languages and the setting is not defined, the default language is `en` if an enabled language named `en` exists, otherwise the first enabled language in the project. The first language is identified by the lowest [_weight_](g), using lexicographical order as the final fallback if weights are tied or undefined.

  See also: [_language_](g).

  [`defaultContentLanguage`]: /configuration/all/#defaultcontentlanguage
