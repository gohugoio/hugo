---
title: Language
description: Returns the language object for the given site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: langs.Language
    signatures: [SITE.Language]
---

The `Language` method on a `Page` object returns the language object for the given page. The language object is derived from the language definition in the site configuration.

You can also use the `Language` method on a `Page` object. See&nbsp;[details].

## Methods

The examples below assume the following in the site configuration:

{{< code-toggle file=hugo >}}
[languages.de]
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
weight = 1
{{< /code-toggle >}}

### IsDefault

{{< new-in 0.153.0 />}}

(`bool`) Reports whether this is the default language object as defined by the [`defaultContentLanguage`][] setting in the site configuration.

```go-html-template
{{ .Site.Language.IsDefault }} → true
```

### Lang

(`string`) The language tag as defined by [RFC 5646]. This is the lower cased key from the site configuration.

```go-html-template
{{ .Site.Language.Lang }} → de
```

### LanguageCode

(`string`) The language code from the site configuration. Falls back to `Lang` if not defined.

```go-html-template
{{ .Site.Language.LanguageCode }} → de-DE
```

### LanguageDirection

(`string`) The language direction from the site configuration, either `ltr` or `rtl`.

```go-html-template
{{ .Site.Language.LanguageDirection }} → ltr
```

### LanguageName

(`string`) The language name from the site configuration.

```go-html-template
{{ .Site.Language.LanguageName }} → Deutsch
```

### Name

(`string`) The language tag as defined by [RFC 5646]. This is the lower cased key from the site configuration. This is an alias for `Lang`.

```go-html-template
{{ .Site.Language.Name }} → de
```

### Weight

(`int`) The language weight from the site configuration which determines its order in the slice of languages returned by the `Languages` method on a `Site` object.

```go-html-template
{{ .Site.Language.Weight }} → 1
```

## Example

Some of the methods above are commonly used in a base template as attributes for the `html` element.

```go-html-template
<html
  lang="{{ .Site.Language.LanguageCode }}" 
  dir="{{ or .Site.Language.LanguageDirection `ltr` }}"
>
```

[`defaultContentLanguage`]: /configuration/all/#defaultcontentlanguage
[details]: /methods/page/language/
[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646
