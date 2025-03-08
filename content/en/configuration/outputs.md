---
title: Configure outputs
linkTitle: Outputs
description: Configure which output formats to render for each page kind.
categories: []
keywords: []
---

{{% glossary-term "output format" %}}

Learn more about creating and configuring output formats in the [configure output formats] section.

## Outputs per page kind

The following default configuration determines the output formats generated for each page kind:

{{< code-toggle config=outputs />}}

To render the built-in `json` output format for the `home` page kind, assuming you've already created the necessary template, add the following to your configuration:

{{< code-toggle file=hugo >}}
[outputs]
home = ['html','rss','json']
{{< /code-toggle >}}

Notice in this example that we only specified the `home` page kind. You don't need to include entries for other page kinds unless you intend to modify their default output formats.

> [!note]
> The order of the output formats in the arrays above is important. The first element will be the _primary output format_ for that page kind, and in most cases that should be `html` as shown in the default configuration.
>
> The primary output format for a given page kind determines the value returned by the [`Permalink`] and [`RelPermalink`] methods on a `Page` object.
>
> See the [link to output formats] section for details.

## Outputs per page

Add output formats to a page's rendering using the `outputs` field in its front matter. For example, to include `json` in the output formats rendered for a specific page:

{{< code-toggle file=content/example.md fm=true >}}
title = 'Example'
outputs = ['json']
{{< /code-toggle >}}

In its default configuration, Hugo will render both the `html` and `json` output formats for this page. The `outputs` field appends to, rather than replaces, the site's configured outputs.

[`Permalink`]: /methods/page/permalink/
[`RelPermalink`]: /methods/page/relpermalink/
[configure output formats]: /configuration/output-formats/
[link to output formats]: configuration/output-formats/#link-to-output-formats
