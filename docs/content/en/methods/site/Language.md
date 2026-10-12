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

You can also use the `Language` method on a `Page` object. See [details][].

## Methods

Use these methods on the `Language` object.

The examples below assume the following language definition.

{{< code-toggle file=hugo >}}
[languages.de]
direction = 'ltr'
label = 'Deutsch'
locale = 'de-DE'
weight = 2
{{< /code-toggle >}}

`Direction`
: {{< new-in 0.158.0 />}}
: (`string`) Returns the [`direction`][] from the language definition.

  ```go-html-template
  {{ .Site.Language.Direction }} → ltr
  ```

`FormatAccounting`
: {{< new-in 0.168.0 />}}
: (`string`) Formats a number as currency in accounting notation for this language. See [`lang.FormatAccounting`][].

  ```go-html-template
  {{ .Site.Language.FormatAccounting 2 "EUR" -12.3 }} → -12,30 €
  ```

`FormatCurrency`
: {{< new-in 0.168.0 />}}
: (`string`) Formats a number as currency for this language. See [`lang.FormatCurrency`][].

  ```go-html-template
  {{ .Site.Language.FormatCurrency 2 "EUR" 12.3 }} → 12,30 €
  ```

`FormatNumber`
: {{< new-in 0.168.0 />}}
: (`string`) Formats a number for this language. See [`lang.FormatNumber`][].

  ```go-html-template
  {{ .Site.Language.FormatNumber 2 1234.5 }} → 1.234,50
  ```

`FormatPercent`
: {{< new-in 0.168.0 />}}
: (`string`) Formats a number as a percentage for this language. See [`lang.FormatPercent`][].

  ```go-html-template
  {{ .Site.Language.FormatPercent 1 12.3 }} → 12,3 %
  ```

`IsDefault`
: {{< new-in 0.153.0 />}}
: (`bool`) Reports whether this is the [default language](g).

  ```go-html-template
  {{ .Site.Language.IsDefault }} → true
  ```

`Label`
: {{< new-in 0.158.0 />}}
: (`string`) Returns the [`label`][] from the language definition.

  ```go-html-template
  {{ .Site.Language.Label }} → Deutsch
  ```

`Lang`
: {{<deprecated-in 0.158.0 />}}
: Use [`Name`](#name) instead.

`LanguageCode`
: {{<deprecated-in 0.158.0 />}}
: Use [`Locale`](#locale) instead.

`LanguageDirection`
: {{<deprecated-in 0.158.0 />}}
: Use [`Direction`](#direction) instead.

`LanguageName`
: {{<deprecated-in 0.158.0 />}}
: Use [`Label`](#label) instead.

`Locale`
: {{< new-in 0.158.0 />}}
: (`string`) Returns the [`locale`][] from the language definition, falling back to [`Name`](#name).

  ```go-html-template
  {{ .Site.Language.Locale }} → de-DE
  ```

`Name`
: {{< new-in 0.153.0 />}}
: (`string`) Returns the language tag as defined by [RFC 5646][]. This is the lowercased key from the language definition.

  ```go-html-template
  {{ .Site.Language.Name }} → de
  ```

`Translate`
: {{< new-in 0.168.0 />}}
: (`string`) Translates a string using the translation tables for this language. See [`lang.Translate`][].

  ```go-html-template
  {{ .Site.Language.Translate "hello" }} → Hallo
  ```

  Use this to translate a string into a language other than that of the current page:

  ```go-html-template
  {{ range .Site.Languages }}
    {{ .Translate "hello" }}
  {{ end }}
  ```

`Weight`
: {{<deprecated-in 0.158.0 />}}

## Example

Some of the methods above are commonly used in a _base_ template as attributes for the `html` element.

```go-html-template
<html
  lang="{{ .Site.Language.Locale }}"
  dir="{{ or .Site.Language.Direction `ltr` }}"
>
```

[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646
[`direction`]: /configuration/languages/#direction
[`lang.FormatAccounting`]: /functions/lang/formataccounting/
[`lang.FormatCurrency`]: /functions/lang/formatcurrency/
[`lang.FormatNumber`]: /functions/lang/formatnumber/
[`lang.FormatPercent`]: /functions/lang/formatpercent/
[`lang.Translate`]: /functions/lang/translate/
[`label`]: /configuration/languages/#label
[`locale`]: /configuration/languages/#locale
[details]: /methods/page/language/
