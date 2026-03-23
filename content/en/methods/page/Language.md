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
{{ .Language.Direction }} → ltr
```

### IsDefault

{{< new-in 0.153.0 />}}

(`bool`) Reports whether this is the [default language][].

```go-html-template
{{ .Language.IsDefault }} → true
```

### Label

{{< new-in 0.158.0 />}}

(`string`) Returns the [`label`][] from the language definition.

```go-html-template
{{ .Language.Label }} → Deutsch
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
{{ .Language.Locale }} → de-DE
```

### Name

{{< new-in 0.153.0 />}}

(`string`) Returns the language tag as defined by [RFC 5646][]. This is the lowercased key from the language definition.

```go-html-template
{{ .Language.Name }} → de
```

### Weight

{{<deprecated-in 0.158.0 />}}

[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646
[`direction`]: /configuration/languages/#direction
[`label`]: /configuration/languages/#label
[`locale`]: /configuration/languages/#locale
[default language]: /quick-reference/glossary/#default-language
[details]: /methods/site/language/

## Example

Use the code below to create a language selector, allowing users to navigate between the different translated versions of the current page.

```go-html-template {file="layouts/_partials/language-selector.html" copy=true}
{{ with .Rotate "language" }}
  <nav class="language-selector">
    <ul>
      {{ range . }}
        {{ if eq .Language $.Language }}
          <li class="active">
            <a aria-current="page" href="{{ .Permalink }}" hreflang="{{ .Language.Locale }}">{{ .Language.Label }}</a>
          </li>
        {{ else }}
          <li>
            <a href="{{ .Permalink }}" hreflang="{{ .Language.Locale }}">{{ .Language.Label }}</a>
          </li>
        {{ end }}
      {{ end }}
    </ul>
  </nav>
{{ end }}
```
