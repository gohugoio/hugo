---
title: Configure params
linkTitle: Params
description: Create custom site parameters.
categories: []
keywords: []
---

Use the `params` key for custom parameters:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/'
title = 'Project Documentation'
languageCode = 'en-US'
[params]
subtitle = 'Reference, Tutorials, and Explanations'
[params.contact]
email = 'info@example.org'
phone = '+1 206-555-1212'
{{< /code-toggle >}}

Access the custom parameters from your templates using the [`Params`] method on a `Site` object:

[`Params`]: /methods/site/params/

```go-html-template
{{ .Site.Params.subtitle }} → Reference, Tutorials, and Explanations
{{ .Site.Params.contact.email }} → info@example.org
```

Key names should use camelCase or snake_case. While TOML, YAML, and JSON allow kebab-case keys, they are not valid [identifiers](g) and cannot be used when [chaining](g) identifiers.

For example, you can do either of these:

```go-html-template
{{ .Site.params.camelCase.foo }}
{{ .Site.params.snake_case.foo }}
```

But you cannot do this:

```go-html-template
{{ .Site.params.kebab-case.foo }}
```

## Multilingual sites

For multilingual sites, create a `params` key under each language:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/'
defaultContentLanguage = 'en'

[languages.de]
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
title = 'Projekt Dokumentation'
weight = 1

[languages.de.params]
subtitle = 'Referenz, Tutorials und Erklärungen'

[languages.de.params.contact]
email = 'info@de.example.org'
phone = '+49 30 1234567'

[languages.en]
languageCode = 'en-US'
languageDirection = 'ltr'
languageName = 'English'
title = 'Project Documentation'
weight = 2

[languages.en.params]
subtitle = 'Reference, Tutorials, and Explanations'

[languages.en.params.contact]
email = 'info@example.org'
phone = '+1 206-555-1212'
{{< /code-toggle >}}

## Namespacing

To prevent naming conflicts, module and theme developers should namespace any custom parameters specific to their module or theme.

{{< code-toggle file=hugo >}}
[params.modules.myModule.colors]
background = '#efefef'
font = '#222222'
{{< /code-toggle >}}

To access the module/theme settings:

```go-html-template
{{ $cfg := .Site.Params.module.mymodule }}

{{ $cfg.colors.background }} → #efefef
{{ $cfg.colors.font }} → #222222
```
