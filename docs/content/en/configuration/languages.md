---
title: Configure languages
linkTitle: Languages
description: Configure the languages in your multilingual project.
categories: []
keywords: []
---

## Base settings

Configure the following base settings:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = false
disableDefaultLanguageRedirect = false
disableLanguages = []
{{< /code-toggle >}}

defaultContentLanguage
: (`string`) The projects's default content language, conforming to the syntax described in [RFC 5646][]. This value must match one of the defined [language keys][]. Default is `en`.

defaultContentLanguageInSubdir
: (`bool`) Whether to publish the default content language to a subdirectory matching the [`defaultContentLanguage`][]. Default is `false`.

disableDefaultLanguageRedirect
: {{< new-in 0.140.0 />}}
: (`bool`) Whether to disable generation of the alias redirect for the default content language. When [`defaultContentLanguageInSubdir`][] is `true`, this setting prevents the root directory from redirecting to the language subdirectory. Conversely, when `defaultContentLanguageInSubdir` is `false`, this setting prevents the language subdirectory from redirecting to the root directory. This is superseded by the more general [`disableDefaultSiteRedirect`][] setting. Default is `false`.

disableLanguages
: (`[]string]`) A slice of language keys representing the languages to disable during the build process. Although this is functional, consider using the [`disabled`](#disabled) key under each language instead.

## Language settings

Configure each language under the `languages` key:

{{< code-toggle config=languages />}}

In the above, `en` is the [language key](#language-keys).

direction
: (`string`) The language direction, either left-to-right (`ltr`) or right-to-left (`rtl`). Use this value in your templates with the global [`dir`][] HTML attribute. Access this value from a template using the [`Language.Direction`][] method on a `Site` or `Page` object. Default is `ltr`.

disabled
: (`bool`) Whether to disable this language when building the site. Default is `false`.

label
: (`string`) The language name, typically used when rendering a language switcher. Access this value from a template using the [`Language.Label`][] method on a `Site` or `Page` object.

languageCode
: {{<deprecated-in 0.158.0 />}}
: Use [`locale`](#locale) instead.

languageDirection
: {{<deprecated-in 0.158.0 />}}
: Use [`direction`](#direction) instead.

languageName
: {{<deprecated-in 0.158.0 />}}
: Use [`label`](#label) instead.

locale
: (`string`) The language tag as described in [RFC 5646][]. This is the primary value used by the [`language.Translate`][] function to select a translation table, falling back to the language key if a matching translation table does not exist.

  Hugo also uses this value to populate:

  - The `lang` attribute of the `html` element in the [embedded alias template][]
  - The `language` element in the [embedded RSS template][]
  - The `locale` property in the [embedded OpenGraph template][]

  > [!note]
  > This value does not affect localization of dates, numbers, and currencies, nor does it affect the site's URL structure. These are controlled by the [language key](#language-keys).

  Access this value from a template using the [`Language.Locale`][] method on a `Site` or `Page` object.

title
: (`string`) The site title for this language. Access this value from a template using the [`Title`][] method on a `Site` object.

weight
: (`int`) The language [weight](g). When set to a non-zero value, this is the primary sort criteria for this language.

## Sort order

Hugo sorts languages by weight in ascending order, then lexicographically in ascending order. This affects build order and complement selection.

## Localized settings

Some configuration settings can be defined separately for each language. For example:

{{< code-toggle file=hugo >}}
[languages.en]
label = 'English'
locale = 'en-US'
timeZone = 'America/New_York'
title = 'Project Documentation'
weight = 1
[languages.en.pagination]
path = 'page'
[languages.en.params]
subtitle = 'Reference, Tutorials, and Explanations'
{{< /code-toggle >}}

The following configuration keys can be defined separately for each language:

{{< per-lang-config-keys >}}

Any key not defined in a `languages` object will fall back to the global value in the root of your project configuration.

## Language keys

Language keys must conform to the syntax described in [RFC 5646][]. For example:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
[languages.de]
weight = 1
[languages.en-US]
weight = 2
[languages.pt-BR]
weight = 3
{{< /code-toggle >}}

Artificial languages with private use subtags as defined in [RFC 5646 § 2.2.7][] are also supported. Omit the `art-x-` prefix from the language key. For example:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'
[languages.en]
weight = 1
[languages.hugolang]
weight = 2
{{< /code-toggle >}}

> [!note]
> Private use subtags must not exceed 8 alphanumeric characters.

## Example

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true
disableDefaultLanguageRedirect = false

[languages.de]
contentDir = 'content/de'
direction = 'ltr'
disabled = false
label = 'Deutsch'
locale = 'de-DE'
title = 'Projekt Dokumentation'
weight = 1

[languages.de.params]
subtitle = 'Referenz, Tutorials und Erklärungen'

[languages.en]
contentDir = 'content/en'
direction = 'ltr'
disabled = false
label = 'English'
locale = 'en-US'
title = 'Project Documentation'
weight = 2

[languages.en.params]
subtitle = 'Reference, Tutorials, and Explanations'
{{< /code-toggle >}}

> [!note]
> In the example above, omit `contentDir` if [translating by file name][].

## Multihost

Hugo supports multiple languages in a multihost configuration. This means you can configure a `baseURL` per `language`.

> [!note]
> If you define a `baseURL` for one language, you must define a unique `baseURL` for all languages.

For example:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'fr'
[languages.en]
baseURL = 'https://en.example.org/'
label = 'English'
title = 'In English'
weight = 2
[languages.fr]
baseURL = 'https://fr.example.org'
label = 'Français'
title = 'En Français'
weight = 1
{{</ code-toggle >}}

With the above, Hugo publishes two sites, each with their own root:

```text
public
├── en
└── fr
```

[RFC 5646 § 2.2.7]: https://datatracker.ietf.org/doc/html/rfc5646#section-2.2.7
[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646#section-2.1
[`Language.Direction`]: /methods/site/language/#direction
[`Language.Label`]: /methods/site/language/#label
[`Language.Locale`]: /methods/site/language/#locale
[`Title`]: /methods/site/title/
[`defaultContentLanguageInSubdir`]: #defaultcontentlanguageinsubdir
[`defaultContentLanguage`]: #defaultcontentlanguage
[`dir`]: https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/dir
[`disableDefaultSiteRedirect`]: /configuration/all/#disabledefaultsiteredirect
[`language.Translate`]: /functions/lang/translate/
[embedded OpenGraph template]: <{{% eturl opengraph %}}>
[embedded RSS template]: <{{% eturl rss %}}>
[embedded alias template]: <{{% eturl alias %}}>
[language keys]: #language-keys
[translating by file name]: /content-management/multilingual/#translation-by-file-name
