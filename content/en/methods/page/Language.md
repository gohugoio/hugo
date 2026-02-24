---
title: Language
description: Returns the Language object for the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: langs.Language
    signatures: [PAGE.Language]
---

The `Language` method on a `Page` object returns the `Language` object for the given page, derived from the language definition in your project configuration.

You can also use the `Language` method on a `Site` object. See&nbsp;[details][].

## Methods

The examples below assume the following in your project configuration:

{{< code-toggle file=hugo >}}
[languages.de]
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
weight = 2
{{< /code-toggle >}}

### IsDefault

{{< new-in 0.153.0 />}}

(`bool`) Reports whether this is the [default language][].

```go-html-template
{{ .Language.IsDefault }} → true
```

### Lang

(`string`) Returns the language tag as defined by [RFC 5646][]. This is the lowercased key from your project configuration.

```go-html-template
{{ .Language.Lang }} → de
```

### LanguageCode

(`string`) Returns the [`languageCode`][] from your project configuration. Falls back to `Lang` if not defined.

```go-html-template
{{ .Language.LanguageCode }} → de-DE
```

### LanguageDirection

(`string`) Returns the [`languageDirection`][] from your project configuration.

```go-html-template
{{ .Language.LanguageDirection }} → ltr
```

### LanguageName

(`string`) Returns the [`languageName`][] from your project configuration.

```go-html-template
{{ .Language.LanguageName }} → Deutsch
```

### Name

{{< new-in 0.153.0 />}}

(`string`) Returns the language tag as defined by [RFC 5646][]. This is the lowercased key from your project configuration. This is an alias for `Lang`.

```go-html-template
{{ .Language.Name }} → de
```

### Weight

(`int`) Returns the language [`weight`][] from your project configuration.

```go-html-template
{{ .Language.Weight }} → 2
```

[`languageCode`]: /configuration/languages/#languagecode
[`languageDirection`]: /configuration/languages/#languagedirection
[`languageName`]: /configuration/languages/#languagename
[`weight`]: /configuration/languages/#weight
[default language]: /quick-reference/glossary/#default-language
[details]: /methods/page/language/
[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646
