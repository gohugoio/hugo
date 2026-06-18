---
_comment: Do not remove front matter.
---

`locale`
: (`string`) The language tag as described in [RFC 5646][]. This is the primary value used by the [`language.Translate`][] function to select a translation table, and for localization of dates, currencies, numbers, and percentages, falling back to the [language key][] in both cases.

  Hugo also uses this value to populate:

  - The `lang` attribute of the `html` element in the [embedded alias template][]
  - The `language` element in the [embedded RSS template][]
  - The `locale` property in the [embedded Open Graph template][]

  Access this value from a template using the [`Language.Locale`][] method on a `Site` or `Page` object.

  [RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646#section-2.1
  [`Language.Locale`]: /methods/site/language/#locale
  [`language.Translate`]: /functions/lang/translate/
  [embedded Open Graph template]: <{{% eturl opengraph %}}>
  [embedded RSS template]: <{{% eturl rss %}}>
  [embedded alias template]: <{{% eturl alias %}}>
  [language key]: /configuration/languages/#language-keys
