---
title: Language
description: Returns the language object for the given page.
categories: []
keywords: []
action:
  related:
    - methods/site/Language
  returnType: langs.Language
  signatures: [PAGE.Language]
---

The `Language` method on a `Page` object returns the language object for the given page. The language object points to the language definition in the site configuration.

You can also use the `Language` method on a `Site` object. See&nbsp;[details].

## Methods

The examples below assume the following in your site configuration:

{{< code-toggle file=hugo >}}
[languages.de]
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
weight = 2
{{< /code-toggle >}}

Lang
: (`string`) The language tag as defined by [RFC 5646].

```go-html-template
{{ .Language.Lang }} → de
```

LanguageCode
: (`string`) The language code from the site configuration. Falls back to `Lang` if not defined.

```go-html-template
{{ .Language.LanguageCode }} → de-DE
```

LanguageDirection
: (`string`) The language direction from the site configuration, either `ltr` or `rtl`.

```go-html-template
{{ .Language.LanguageDirection }} → ltr
```

LanguageName
: (`string`) The language name from the site configuration.

```go-html-template
{{ .Language.LanguageName }} → Deutsch
```

Weight
: (`int`) The language weight from the site configuration which determines its order in the slice of languages returned by the `Languages` method on a `Site` object.

```go-html-template
{{ .Language.Weight }} → 2
```

[details]: /methods/site/language/
[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646
