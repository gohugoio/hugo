---
title: Site variables
description: Use these methods with Site objects. A multilingual project will have two or more sites, one for each language.
categories: [variables]
keywords: [site]
menu:
  docs:
    parent: variables
    weight: 80
weight: 80
aliases: [/variables/site-variables/]
toc: true
---

{{% include "variables/_common/consistent-terminology.md" %}}

## All methods

Use any of these methods in your templates.

{{< list-pages-in-section path=/methods/site titlePrefix=.Site. >}}

## Multilingual

Use these methods with your multilingual projects.

{{< list-pages-in-section path=/methods/site filter=methods_site_multilingual filterType=include titlePrefix=.Site. omitElementIDs=true >}}

[`site`]: /functions/global/site
[context]: /getting-started/glossary/#context
[configuration file]: /getting-started/configuration

## Page collections

Range through these collections when rendering lists on any page.

{{< list-pages-in-section path=/methods/site filter=methods_site_page_collections filterType=include titlePrefix=.Site. omitElementIDs=true >}}

## Global site function

Within a partial template, if you did not pass a `Page` or `Site` object in [context], you cannot use this syntax:

```go-html-template
{{ .Site.SomeMethod }}
```

Instead, use the global [`site`] function:

```go-html-template
{{ site.SomeMethod }}
```

{{% note %}}
You can use the global site function in all templates to avoid context problems. Its usage is not limited to partial templates.
{{% /note %}}
