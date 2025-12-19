---
title: Configure languages
linkTitle: Languages
description: Configure the languages in your multilingual site.
categories: []
keywords: []
---

## Base settings

Configure the following base settings within the site's root configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = false
disableDefaultLanguageRedirect = false
disableLanguages = []
{{< /code-toggle >}}

defaultContentLanguage
: (`string`) The project's default language key, conforming to the syntax described in [RFC 5646]. This value must match one of the defined [language keys](#language-keys). Default is `en`.

defaultContentLanguageInSubdir
: (`bool`) Whether to publish the default language site to a subdirectory matching the `defaultContentLanguage`. Default is `false`.

disableDefaultLanguageRedirect
: {{< new-in 0.140.0 />}}
: (`bool`) Whether to disable generation of the alias redirect to the default language when `DefaultContentLanguageInSubdir` is `true`. Default is `false`.

disableLanguages
: (`[]string]`) A slice of language keys representing the languages to disable during the build process. Although this is functional, consider using the [`disabled`](#disabled) key under each language instead.

## Language settings

Configure each language under the `languages` key:

{{< code-toggle config=languages />}}

In the above, `en` is the [language key](#language-keys).

disabled
: (`bool`) Whether to disable this language when building the site. Default is `false`.

languageCode
: (`string`) The language tag as described in [RFC 5646]. This value does not affect localization or URLs. Hugo uses this value to populate:

  - The `lang` attribute of the `html` element in the [embedded alias template]
  - The `language` element in the [embedded RSS template]
  - The `locale` property in the [embedded OpenGraph template]

  Access this value from a template using the [`Language.LanguageCode`] method on a `Site` or `Page` object.

languageDirection
: (`string`) The language direction, either left-to-right (`ltr`) or right-to-left (`rtl`). Use this value in your templates with the global [`dir`] HTML attribute. Access this value from a template using the [`Language.LanguageDirection`] method on a `Site` or `Page` object.

languageName
: (`string`) The language name, typically used when rendering a language switcher. Access this value from a template using the [`Language.LanguageName`] method on a `Site` or `Page` object.

title
: (`string`) The site title for this language. Access this value from a template using the [`Title`] method on a `Site` object.

weight
: (`int`) The language [weight](g). When set to a non-zero value, this is the primary sort criteria for this language. Access this value from a template using the [`Language.Weight`] method on a `Site` or `Page` object.

## Localized settings

Some configuration settings can be defined separately for each language. For example:

{{< code-toggle file=hugo >}}
[languages.en]
languageCode = 'en-US'
languageName = 'English'
weight = 1
title = 'Project Documentation'
timeZone = 'America/New_York'
[languages.en.pagination]
path = 'page'
[languages.en.params]
subtitle = 'Reference, Tutorials, and Explanations'
{{< /code-toggle >}}

The following configuration keys can be defined separately for each language:

{{< per-lang-config-keys >}}

Any key not defined in a `languages` object will fall back to the global value in the root of the site configuration.

## Language keys

Language keys must conform to the syntax described in [RFC 5646]. For example:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
[languages.de]
  weight = 1
[languages.en-US]
  weight = 2
[languages.pt-BR]
  weight = 3
{{< /code-toggle >}}

Artificial languages with private use subtags as defined in [RFC 5646 § 2.2.7] are also supported. Omit the `art-x-` prefix from the language key. For example:

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
disabled = false
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
title = 'Projekt Dokumentation'
weight = 1

[languages.de.params]
subtitle = 'Referenz, Tutorials und Erklärungen'

[languages.en]
contentDir = 'content/en'
disabled = false
languageCode = 'en-US'
languageDirection = 'ltr'
languageName = 'English'
title = 'Project Documentation'
weight = 2

[languages.en.params]
subtitle = 'Reference, Tutorials, and Explanations'
{{< /code-toggle >}}

> [!note]
> In the example above, omit `contentDir` if [translating by file name].

## Multihost

Hugo supports multiple languages in a multihost configuration. This means you can configure a `baseURL` per `language`.

> [!note]
> If you define a `baseURL` for one language, you must define a unique `baseURL` for all languages.

For example:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'fr'
[languages]
  [languages.en]
    baseURL = 'https://en.example.org/'
    languageName = 'English'
    title = 'In English'
    weight = 2
  [languages.fr]
    baseURL = 'https://fr.example.org'
    languageName = 'Français'
    title = 'En Français'
    weight = 1
{{</ code-toggle >}}

With the above, Hugo publishes two sites, each with their own root:

```text
public
├── en
└── fr
```

[`dir`]: https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/dir
[`Language.LanguageCode`]: /methods/site/language/#languagecode
[`Language.LanguageDirection`]: /methods/site/language/#languagedirection
[`Language.LanguageName`]: /methods/site/language/#languagename
[`Language.Weight`]: /methods/site/language/#weight
[`Title`]: /methods/site/title/
[embedded alias template]: <{{% eturl alias %}}>
[embedded OpenGraph template]: <{{% eturl opengraph %}}>
[embedded RSS template]: <{{% eturl rss %}}>
[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646#section-2.1
[RFC 5646 § 2.2.7]: https://datatracker.ietf.org/doc/html/rfc5646#section-2.2.7
[translating by file name]: /content-management/multilingual/#translation-by-file-name
