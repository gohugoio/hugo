---
title: Params
description: Returns a map of custom parameters as defined in the site configuration.
categories: []
keywords: []
action:
  related:
    - functions/collections/indexFunction
    - methods/page/Params
    - methods/page/Param
  returnType: maps.Params
  signatures: [SITE.Params]
---

With this site configuration:

{{< code-toggle file=hugo >}}
[params]
  subtitle = 'The Best Widgets on Earth'
  copyright-year = '2023'
  [params.author]
    email = 'jsmith@example.org'
    name = 'John Smith'
  [params.layouts]
    rfc_1123 = 'Mon, 02 Jan 2006 15:04:05 MST'
    rfc_3339 = '2006-01-02T15:04:05-07:00'
{{< /code-toggle >}}

Access the custom parameters by [chaining] the [identifiers]:

```go-html-template
{{ .Site.Params.subtitle }} → The Best Widgets on Earth
{{ .Site.Params.author.name }} → John Smith

{{ $layout := .Site.Params.layouts.rfc_1123 }}
{{ .Site.Lastmod.Format $layout }} → Tue, 17 Oct 2023 13:21:02 PDT
```

In the template example above, each of the keys is a valid identifier. For example, none of the keys contains a hyphen. To access a key that is not a valid identifier, use the [`index`] function:

```go-html-template
{{ index .Site.Params "copyright-year" }} → 2023
```

[`index`]: /functions/collections/indexfunction/
[chaining]: /getting-started/glossary/#chain
[identifiers]: /getting-started/glossary/#identifier
