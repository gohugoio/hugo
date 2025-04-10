---
title: Configure media types
linkTitle: Media types
description: Configure media types.
categories: []
keywords: []
---

{{% glossary-term "media type" %}}

Configured media types serve multiple purposes in Hugo, including the definition of [output formats](g). This is the default media type configuration in tabular form:

{{< datatable "config" "mediaTypes" "_key" "suffixes" >}}

The `suffixes` column in the table above shows the suffixes associated with each media type. For example, Hugo associates `.html` and `.htm` files with the `text/html` media type.

> [!note]
> The first suffix is the primary suffix. Use the primary suffix when naming template files. For example, when creating a template for an RSS feed, use the `xml` suffix.

## Default configuration

The following is the default configuration that matches the table above:

{{< code-toggle file=hugo config=mediaTypes />}}

delimiter
: (`string`) The delimiter between the file name and the suffix. The delimiter, in conjunction with the suffix, forms the file extension. Default is `"."`.

suffixes
: (`[]string`) The suffixes associated with this media type. The first suffix is the primary suffix.

## Modify a media type

You can modify any of the default media types. For example, to switch the primary suffix for `text/html` from `html` to `htm`:

{{< code-toggle file=hugo >}}
[mediaTypes.'text/html']
suffixes = ['htm','html']
{{< /code-toggle >}}

If you alter a default media type, you must also explicitly redefine all output formats that utilize that media type. For example, to ensure the changes above affect the `html` output format, redefine the `html` output format:

{{< code-toggle file=hugo >}}
[outputFormats.html]
mediaType = 'text/html'
{{< /code-toggle >}}

## Create a media type

You can create new media types as needed. For example, to create a media type for an Atom feed:

{{< code-toggle file=hugo >}}
[mediaTypes.'application/atom+xml']
suffixes = ['atom']
{{< /code-toggle >}}

## Media types without suffixes

Occasionally, you may need to create a media type without a suffix or delimiter. For example, [Netlify] recognizes configuration files named `_redirects` and `_headers`, which Hugo can generate using custom [output formats](g).

To support these custom output formats, register a custom media type with no suffix or delimiter:

{{< code-toggle file=hugo >}}
[mediaTypes."text/netlify"]
delimiter = ""
{{< /code-toggle >}}

The custom output format definitions would look something like this:

{{< code-toggle file=hugo >}}
[outputFormats.redir]
baseName    = "_redirects"
isPlainText = true
mediatype   = "text/netlify"
[outputFormats.headers]
baseName       = "_headers"
isPlainText    = true
mediatype      = "text/netlify"
notAlternative = true
{{< /code-toggle >}}

[Netlify]: https://www.netlify.com/
