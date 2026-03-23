---
title: Language
description: Returns the Language object for the given site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: langs.Language
    signatures: [SITE.Language]
---

The `Language` method on a `Site` object returns the `Language` object for the given site, derived from the language definition in your project configuration.

You can also use the `Language` method on a `Page` object. See&nbsp;[details][].

## Methods

The examples below assume the following language definition.

{{< code-toggle file=hugo >}}
[languages.de]
direction = 'ltr'
label = 'Deutsch'
locale = 'de-DE'
weight = 2
{{< /code-toggle >}}

### Direction

{{< new-in 0.158.0 />}}

(`string`) Returns the [`direction`][] from the language definition.

```go-html-template
{{ .Site.Language.Direction }} → ltr
```

### IsDefault

{{< new-in 0.153.0 />}}

(`bool`) Reports whether this is the [default language][].

```go-html-template
{{ .Site.Language.IsDefault }} → true
```

### Label

{{< new-in 0.158.0 />}}

(`string`) Returns the [`label`][] from the language definition.

```go-html-template
{{ .Site.Language.Label }} → Deutsch
```

### Lang

{{<deprecated-in 0.158.0 />}}

Use [`Name`](#name) instead.

### LanguageCode

{{<deprecated-in 0.158.0 />}}

Use [`Locale`](#locale) instead.

### LanguageDirection

{{<deprecated-in 0.158.0 />}}

Use [`Direction`](#direction) instead.

### LanguageName

{{<deprecated-in 0.158.0 />}}

Use [`Label`](#label) instead.

### Locale

{{< new-in 0.158.0 />}}

(`string`) Returns the [`locale`][] from the language definition, falling back to [`Name`](#name).

```go-html-template
{{ .Site.Language.Locale }} → de-DE
```

### Name

{{< new-in 0.153.0 />}}

(`string`) Returns the language tag as defined by [RFC 5646][]. This is the lowercased key from the language definition.

```go-html-template
{{ .Site.Language.Name }} → de
```

### Weight

{{<deprecated-in 0.158.0 />}}

## Example

Some of the methods above are commonly used in a base template as attributes for the `html` element.

```go-html-template
<html
  lang="{{ .Site.Language.Locale }}" 
  dir="{{ or .Site.Language.Direction `ltr` }}"
>
```

[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646
[`direction`]: /configuration/languages/#direction
[`label`]: /configuration/languages/#label
[`locale`]: /configuration/languages/#locale
[default language]: /quick-reference/glossary/#default-language
[details]: /methods/page/language/
